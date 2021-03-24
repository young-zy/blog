package common

import "time"

// IntToUintPointer converts an int to a uint pointer
func IntToUintPointer(x int) *uint {
	y := uint(x)
	return &y
}

// Now returns the pointer type of current time
func Now() *time.Time {
	tmp := time.Now()
	return &tmp
}
