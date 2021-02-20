package common

func IntToUintPointer(x int) *uint {
	y := uint(x)
	return &y
}
