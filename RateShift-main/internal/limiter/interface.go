// Package limiter defines common interface types and data structure
// used by all rate limiting algorithms
package limiter

import (
	"context"
	"time"
)

type AlgorithmType string

const (
	TokenBucketAlgo   AlgorithmType = "token_bucket"
	LeakyBucketAlgo   AlgorithmType = "leaky_bucket"
	SlidingWindowAlgo AlgorithmType = "sliding_window"
)

type Decision struct {
	Allowed    bool
	Remaining  int
	ResetAfter time.Time
}

type State struct {
	Algorithm AlgorithmType
	Data      map[string]any
}

type RateLimiter interface {
	Allow(ctx context.Context, tenantID string) (Decision, error)
	ExportState(ctx context.Context, tenantID string) (State, error)
	ImportState(ctx context.Context, tenantID string, state State) error
	Type() AlgorithmType
}
