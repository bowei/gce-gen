/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloud

import (
	"context"

	"github.com/bowei/gce-gen/pkg/cloud/meta"
)

// RateLimitKey is a key identifying the operation to be rate limited. The rate limit
// queue will be determined based on the contents of RateKey.
type RateLimitKey struct {
	ProjectID string
	Operation string
	Version   meta.Version
	Service   string
}

// RateLimiter is the interface for a rate limiting policy.
type RateLimiter interface {
	Accept(ctx context.Context, key *RateLimitKey)
}

// NopRateLimiter is a rate limiter that performs no limiting.
type NopRateLimiter struct {
}

func (*NopRateLimiter) Accept(ctx context.Context, key *RateLimitKey) {
}
