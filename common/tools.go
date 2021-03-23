package common

import "time"

func IntToUintPointer(x int) *uint {
	y := uint(x)
	return &y
}

func Now() *time.Time {
	tmp := time.Now()
	return &tmp
}
