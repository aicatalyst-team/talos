// Package ratelimit provides rate limit enforcement for API key verification.
// Commercial-only: OSS builds use a no-op implementation.
//
// The Limiter interface, Result, and NoopLimiter types are defined in
// internal/ratelimit and re-exported here for backward compatibility
// within the commercial package tree.
package ratelimit

import (
	internalrl "github.com/ory/talos/internal/ratelimit"
)

// Limiter is re-exported from internal/ratelimit.
type Limiter = internalrl.Limiter

// Result is re-exported from internal/ratelimit.
type Result = internalrl.Result

// NoopLimiter is re-exported from internal/ratelimit.
type NoopLimiter = internalrl.NoopLimiter
