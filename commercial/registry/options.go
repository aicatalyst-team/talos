// Package registry provides proprietary persistence options and callbacks.
//
// The shared types (FeatureOptions, CacheFactory, RateLimiterFactory, etc.)
// are defined in internal/registrytypes and re-exported here for backward
// compatibility within the commercial package tree.
package registry

import (
	"github.com/ory/talos/internal/registrytypes"
)

// RegisterDatabaseMetricsFunc is re-exported from internal/registrytypes.
type RegisterDatabaseMetricsFunc = registrytypes.RegisterDatabaseMetricsFunc

// HTTPMiddlewareFunc is re-exported from internal/registrytypes.
type HTTPMiddlewareFunc = registrytypes.HTTPMiddlewareFunc

// CacheFactory is re-exported from internal/registrytypes.
type CacheFactory = registrytypes.CacheFactory

// RateLimiterFactory is re-exported from internal/registrytypes.
type RateLimiterFactory = registrytypes.RateLimiterFactory

// FeatureOptions is re-exported from internal/registrytypes.
type FeatureOptions = registrytypes.FeatureOptions
