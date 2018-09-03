package gxlog

type gatherer struct {
	seq int
}

func (this *gatherer) gather() int {
	this.seq++
	return this.seq
}
