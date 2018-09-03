package gxlog

type gatherer struct{}

func (this *gatherer) gather() *Record {
	return &Record{}
}
