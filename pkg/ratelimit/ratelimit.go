// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

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
