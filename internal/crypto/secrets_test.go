package crypto_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/talos/internal/crypto"
	"github.com/ory/talos/internal/testutil"

	"github.com/ory/x/configx"
)

func TestHMACSecretsForVerification(t *testing.T) {
	t.Parallel()

	t.Run("HMAC-specific secret configured with retired", func(t *testing.T) {
		t.Parallel()

		provider := testutil.NewTestProvider(t, configx.WithValues(map[string]any{
			"secrets": map[string]any{
				"hmac": map[string]any{
					"current": "hmac-secret-32chars-minimum-12345678901234",
					"retired": []string{"old-hmac-secret-32chars-1234567890123456"},
				},
			},
		}))
		ctx := context.Background()

		secrets, err := crypto.HMACSecretsForVerification(ctx, provider)
		require.NoError(t, err)
		assert.Equal(t, []string{
			"hmac-secret-32chars-minimum-12345678901234",
			"old-hmac-secret-32chars-1234567890123456",
		}, secrets)
	})

	t.Run("HMAC configured with empty retired", func(t *testing.T) {
		t.Parallel()

		provider := testutil.NewTestProvider(t, configx.WithValues(map[string]any{
			"secrets": map[string]any{
				"hmac": map[string]any{
					"current": "hmac-secret-32chars-minimum-12345678901234",
					"retired": []string{},
				},
			},
		}))
		ctx := context.Background()

		secrets, err := crypto.HMACSecretsForVerification(ctx, provider)
		require.NoError(t, err)
		assert.Equal(t, []string{
			"hmac-secret-32chars-minimum-12345678901234",
		}, secrets)
	})

	// Note: the "HMAC not configured" path is unreachable in production because
	// the config schema requires secrets.hmac.current. Schema-level enforcement
	// is covered by TestNewProvider_RequiresSecretsAtSchemaLevel in
	// internal/config. The defensive runtime guard here remains as belt-and-
	// suspenders.
}

func TestHMACSecretForSigning(t *testing.T) {
	t.Parallel()

	t.Run("HMAC-specific secret configured", func(t *testing.T) {
		t.Parallel()

		provider := testutil.NewTestProvider(t, configx.WithValues(map[string]any{
			"secrets": map[string]any{
				"hmac": map[string]any{
					"current": "hmac-secret-32chars-minimum-12345678901234",
				},
			},
		}))
		ctx := context.Background()

		secret, err := crypto.HMACSecretForSigning(ctx, provider)
		require.NoError(t, err)
		assert.Equal(t, "hmac-secret-32chars-minimum-12345678901234", secret)
	})

	// Note: the "HMAC not configured" path is unreachable in production because
	// the config schema requires secrets.hmac.current.
}

// reviewed - @aeneasr - 2026-03-26
