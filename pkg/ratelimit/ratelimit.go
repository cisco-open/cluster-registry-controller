// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ratelimit

import (
	"emperror.dev/errors"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

var defaultRateQuota = throttled.RateQuota{MaxRate: throttled.PerSec(1), MaxBurst: 0}

func NewRateLimiter(maxKeys int, quota *throttled.RateQuota) (throttled.RateLimiter, error) {
	var rateLimiter *throttled.GCRARateLimiter

	store, err := memstore.New(maxKeys)
	if err != nil {
		return nil, errors.WrapIf(err, "could not create memstore for rate limit")
	}

	if quota == nil {
		quota = &defaultRateQuota
	}

	rateLimiter, err = throttled.NewGCRARateLimiter(store, *quota)
	if err != nil {
		return nil, errors.WrapIf(err, "could not create rate limiter")
	}

	return rateLimiter, nil
}
