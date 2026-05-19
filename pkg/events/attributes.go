package events

import "github.com/ory/x/otelx/semconv"

// Attribute keys used by Talos audit events. These follow OpenTelemetry
// semantic conventions and use PascalCase to match other services in the Ory
// monorepo (Kratos, Hydra, Keto).
const (
	// AttrKeyID identifies the credential involved in an event.
	AttrKeyID semconv.AttributeKey = "APIKeyID"

	// AttrAPIKeyPrefix carries the API key prefix (the non-secret portion).
	AttrAPIKeyPrefix semconv.AttributeKey = "APIKeyPrefix"

	// AttrKeyType distinguishes issued vs imported credentials.
	AttrKeyType semconv.AttributeKey = "KeyType"

	// AttrOperation names the operation that produced the event (for example
	// "create", "revoke").
	AttrOperation semconv.AttributeKey = "Operation"

	// AttrReason is a short human-readable explanation attached to the event.
	AttrReason semconv.AttributeKey = "Reason"

	// AttrActorID identifies the principal that triggered the event.
	AttrActorID semconv.AttributeKey = "ActorID"

	// AttrExpiry is the RFC 3339 timestamp at which the credential expires.
	AttrExpiry semconv.AttributeKey = "Expiry"

	// AttrVisibility describes who can see the credential (for example
	// "public" or "private").
	AttrVisibility semconv.AttributeKey = "Visibility"

	// AttrMetadataPrefix is a string (not AttributeKey) because it is
	// concatenated with user-defined metadata keys rather than used as a
	// standalone attribute key.
	AttrMetadataPrefix = "metadata."
)
