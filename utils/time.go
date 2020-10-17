package utils

import "time"

func Now() uint32 {
	return uint32(time.Now().Unix())
}
