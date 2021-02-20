package common

type Pager struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func NewPager() *Pager {
	return &Pager{
		Page: 1,
		Size: 10,
	}
}
