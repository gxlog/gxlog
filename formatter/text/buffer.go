package text

type buffer struct {
	bufSize  int
	bufCount int

	buf []byte
}

func (this *buffer) Get() []byte {
	if len(this.buf) < this.bufSize {
		this.buf = make([]byte, this.bufSize*this.bufCount)
	}
	buf := this.buf[:0:this.bufSize]
	this.buf = this.buf[this.bufSize:]
	return buf
}

func (this *buffer) GetConfig() (bufSize, bufCount int) {
	return this.bufSize, this.bufCount
}

func (this *buffer) SetConfig(bufSize, bufCount int) {
	this.bufSize = bufSize
	this.bufCount = bufCount
}
