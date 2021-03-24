package common

// Pager is a struct for paging
type Pager struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

// NewPager returns a pager
func NewPager() *Pager {
	return &Pager{
		Page: 1,
		Size: 10,
	}
}
