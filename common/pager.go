package common

type Pager struct {
	Page int
	Size int
}

func NewPager() *Pager {
	return &Pager{
		Page: 1,
		Size: 10,
	}
}
