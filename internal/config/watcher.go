package config

import (
	stderrors "errors"
	"log/slog"

	"github.com/ory/jsonschema/v3"
	"github.com/ory/x/configx"
	"github.com/ory/x/watcherx"
)

// logConfigChange logs the outcome of a configuration hot-reload, mirroring
// configx.LogrusWatcher but using slog (Talos logs with slog only).
//
// configx fires the registered watcher on every config file change. On a
// failed reload (invalid config, or a change to an immutable key) configx rolls
// back to the last working revision; without this log line the rollback is
// silent and operators cannot tell a rejected change from an applied one.
func logConfigChange(l *slog.Logger, e watcherx.Event, err error) {
	src := slog.String("file", e.Source())

	if _, ok := stderrors.AsType[*jsonschema.ValidationError](err); ok {
		l.Error("The changed configuration is invalid and could not be loaded. "+
			"Rolling back to the last working configuration revision. Please address "+
			"the validation errors before restarting the process.",
			src, slog.String("error", err.Error()))
		return
	}

	if immutableErr, ok := stderrors.AsType[*configx.ImmutableError](err); ok {
		// Log only the key, never the values: immutable config keys can hold
		// secrets (DSNs, signing keys, TLS material) that must not leak to logs.
		l.Error("A configuration value marked as immutable has changed. "+
			"Rolling back to the last working configuration revision. To reload the "+
			"values please restart the process.",
			src,
			slog.String("key", immutableErr.Key))
		return
	}

	if err != nil {
		l.Error("An error occurred while watching the configuration file.",
			src, slog.String("error", err.Error()))
		return
	}

	l.Info("Configuration change processed successfully.", src)
}
