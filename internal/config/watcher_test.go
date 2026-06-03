package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/jsonschema/v3"
	"github.com/ory/x/configx"
	"github.com/ory/x/watcherx"
)

func TestLogConfigChange(t *testing.T) {
	const source = "/tmp/talos.yaml"

	for _, tc := range []struct {
		name        string
		err         error
		wantLevel   string
		wantMsg     string
		wantAttrs   map[string]string
		absentAttrs []string
	}{
		{
			name:      "success",
			err:       nil,
			wantLevel: "INFO",
			wantMsg:   "Configuration change processed successfully.",
		},
		{
			name:      "validation error",
			err:       &jsonschema.ValidationError{Message: "value must be a string"},
			wantLevel: "ERROR",
			wantMsg:   "The changed configuration is invalid and could not be loaded. Rolling back to the last working configuration revision. Please address the validation errors before restarting the process.",
		},
		{
			name:      "immutable error",
			err:       configx.NewImmutableError("db.dsn", "old-dsn", "new-dsn"),
			wantLevel: "ERROR",
			wantMsg:   "A configuration value marked as immutable has changed. Rolling back to the last working configuration revision. To reload the values please restart the process.",
			wantAttrs: map[string]string{
				"key": "db.dsn",
			},
			// The values can hold secrets (DSNs, keys) and must never be logged.
			absentAttrs: []string{"old_value", "new_value", "error"},
		},
		{
			name:      "generic error",
			err:       errors.New("disk read failed"),
			wantLevel: "ERROR",
			wantMsg:   "An error occurred while watching the configuration file.",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

			logConfigChange(l, watcherx.NewErrorEvent(nil, source), tc.err)

			var rec map[string]any
			require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))

			assert.Equal(t, tc.wantLevel, rec["level"])
			assert.Equal(t, tc.wantMsg, rec["msg"])
			assert.Equal(t, source, rec["file"])

			for k, want := range tc.wantAttrs {
				assert.Equal(t, want, rec[k], "attr %q", k)
			}

			for _, k := range tc.absentAttrs {
				assert.NotContains(t, rec, k, "attr %q must not be logged", k)
			}

			if tc.err != nil {
				// Validation and generic errors carry the error text; the
				// immutable case reports only the key (values may be secrets).
				if len(tc.wantAttrs) == 0 {
					assert.Equal(t, tc.err.Error(), rec["error"])
				}
			} else {
				assert.NotContains(t, rec, "error")
			}
		})
	}
}
