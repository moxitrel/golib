package svc

import (
	"math"
	"time"
)

var (
	FuncArgMax = math.MaxUint16

	PoolMin     uint32 = 2
	PoolMax     uint32 = math.MaxUint16
	PoolDelay          = 200 * time.Millisecond //a proper value should at least 0.1s
	PoolTimeOut        = time.Minute
)
