package cmd

import (
	"bytes"
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/cmdx"
)

// syncBuffer is a concurrency-safe bytes.Buffer. The background command
// goroutine writes to it while the test goroutine reads from it.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// newTestDSN returns a file-based SQLite DSN in a fresh temp directory. The
// busy timeout makes concurrent writers (e.g. 'migrate up' racing the
// '--block' poll loop's version read) wait instead of failing with
// SQLITE_BUSY.
func newTestDSN(t *testing.T) string {
	t.Helper()
	return "sqlite3://" + filepath.Join(t.TempDir(), "test.db") + "?_pragma=busy_timeout(10000)"
}

func TestMigrateStatusBlock_FullyMigratedReturnsImmediately(t *testing.T) {
	t.Parallel()

	dsn := newTestDSN(t)

	_, upStderr, err := cmdx.ExecCtx(t.Context(), NewRoot(), nil, "migrate", "up", "--database", dsn)
	require.NoError(t, err, "migrate up failed: %s", upStderr)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)

	_, stderr, err := cmdx.ExecCtx(ctx, NewRoot(), nil, "migrate", "status", "--block", "--database", dsn)
	require.NoError(t, err, "status --block on a fully migrated database should return immediately: %s", stderr)
	assert.Contains(t, stderr, "Current Version")
	assert.Contains(t, stderr, "Latest Available")
	assert.NotContains(t, stderr, "Waiting", "should not wait when all migrations are applied")
}

func TestMigrateStatusBlock_UnblocksAfterMigrateUp(t *testing.T) {
	t.Parallel()

	dsn := newTestDSN(t)

	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	t.Cleanup(cancel)

	var stdout, stderr syncBuffer
	eg := cmdx.ExecBackgroundCtx(ctx, NewRoot(), nil, &stdout, &stderr,
		"migrate", "status", "--block", "--database", dsn)

	require.Eventually(t, func() bool {
		return strings.Contains(stderr.String(), "Waiting")
	}, 10*time.Second, 50*time.Millisecond,
		"status --block should report waiting on an uninitialized database; stderr: %s", stderr.String())

	// The command must still be blocked: no final status printed yet.
	assert.NotContains(t, stderr.String(), "Current Version", "should not have printed final status while migrations are pending")

	_, upStderr, err := cmdx.ExecCtx(t.Context(), NewRoot(), nil, "migrate", "up", "--database", dsn)
	require.NoError(t, err, "migrate up failed: %s", upStderr)

	require.NoError(t, eg.Wait(), "status --block should unblock after migrate up; stderr: %s", stderr.String())
	assert.Contains(t, stderr.String(), "Current Version")
	assert.Contains(t, stderr.String(), "Latest Available")
}

func TestMigrateStatusBlock_ContextCancellationUnblocks(t *testing.T) {
	t.Parallel()

	dsn := newTestDSN(t)

	ctx, cancel := context.WithTimeout(t.Context(), 300*time.Millisecond)
	t.Cleanup(cancel)

	start := time.Now()
	_, stderr, err := cmdx.ExecCtx(ctx, NewRoot(), nil, "migrate", "status", "--block", "--database", dsn)
	require.Error(t, err, "status --block must return an error when the context is cancelled while migrations are pending; stderr: %s", stderr)
	assert.Less(t, time.Since(start), 5*time.Second, "cancellation should unblock promptly")
}

func TestMigrateStatus_UninitializedWithoutBlock(t *testing.T) {
	t.Parallel()

	dsn := newTestDSN(t)

	_, stderr, err := cmdx.ExecCtx(t.Context(), NewRoot(), nil, "migrate", "status", "--database", dsn)
	require.NoError(t, err, "status without --block must exit 0 on an uninitialized database: %s", stderr)
	assert.Contains(t, stderr, "Not initialized")
	assert.Contains(t, stderr, "Latest Available")
}

func TestMigrateStatusBlock_DirtyTreatedAsPending(t *testing.T) {
	t.Parallel()

	dbFile := filepath.Join(t.TempDir(), "test.db")
	dsn := "sqlite3://" + dbFile

	_, upStderr, err := cmdx.ExecCtx(t.Context(), NewRoot(), nil, "migrate", "up", "--database", dsn)
	require.NoError(t, err, "migrate up failed: %s", upStderr)

	// Flip the dirty bit directly. Documented exception to the no-direct-DB-writes
	// rule: no production code path can produce a dirty database on demand.
	db, err := sql.Open("sqlite", dbFile)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	_, err = db.ExecContext(t.Context(), "UPDATE schema_migrations SET dirty = 1")
	require.NoError(t, err)
	require.NoError(t, db.Close())

	ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
	t.Cleanup(cancel)

	_, stderr, err := cmdx.ExecCtx(ctx, NewRoot(), nil, "migrate", "status", "--block", "--database", dsn)
	require.Error(t, err, "status --block must keep waiting on a dirty database; stderr: %s", stderr)
	assert.Contains(t, stderr, "Waiting", "dirty database should be treated as pending while blocking")

	_, stderr, err = cmdx.ExecCtx(t.Context(), NewRoot(), nil, "migrate", "status", "--database", dsn)
	require.NoError(t, err, "status without --block must exit 0 on a dirty database: %s", stderr)
	assert.Contains(t, stderr, "DIRTY")
}
