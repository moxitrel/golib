package svc

import (
	"math"
	"time"
)

var (
	POOL_MIN     uint16 = 0
	POOL_MAX     uint16 = math.MaxUint16
	POOL_DELAY          = 200 * time.Millisecond
	POOL_TIMEOUT        = time.Minute
)
