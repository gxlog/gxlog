package gxlog

type Slot int

const (
	Slot0 Slot = iota
	Slot1
	Slot2
	Slot3
	Slot4
	Slot5
	Slot6
	Slot7
	MaxSlot
)

type link struct {
	formatter Formatter
	writer    Writer
	level     Level
	filter    Filter
}

func (this *logger) Link(slot Slot, ft Formatter, wt Writer, opts ...interface{}) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.slots[slot] != nil {
		return false
	}
	this.link(slot, ft, wt, opts)
	return true
}

func (this *logger) ForceLink(slot Slot, ft Formatter, wt Writer, opts ...interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.link(slot, ft, wt, opts)
}

func (this *logger) MustLink(slot Slot, ft Formatter, wt Writer, opts ...interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.slots[slot] != nil {
		panic("slot in use")
	}
	this.link(slot, ft, wt, opts)
}

func (this *logger) Unlink(slot Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot] = nil
}

func (this *logger) UnlinkAll() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for i := range this.slots {
		this.slots[i] = nil
	}
}

func (this *logger) CopyLink(dst, src Slot) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.slots[dst] != nil {
		return false
	}
	this.slots[dst] = this.slots[src]
	return true
}

func (this *logger) ForceCopyLink(dst, src Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[dst] = this.slots[src]
}

func (this *logger) MustCopyLink(dst, src Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.slots[dst] != nil {
		panic("slot in use")
	}
	this.slots[dst] = this.slots[src]
}

func (this *logger) MoveLink(to, from Slot) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.slots[to] != nil {
		return false
	}
	this.slots[to] = this.slots[from]
	this.slots[from] = nil
	return true
}

func (this *logger) ForceMoveLink(to, from Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[to] = this.slots[from]
	this.slots[from] = nil
}

func (this *logger) MustMoveLink(to, from Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.slots[to] != nil {
		panic("slot in use")
	}
	this.slots[to] = this.slots[from]
	this.slots[from] = nil
}

func (this *logger) SwapLink(left, right Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[left], this.slots[right] = this.slots[right], this.slots[left]
}

func (this *logger) HasLink(slot Slot) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot] != nil
}

func (this *logger) SetLinkFormatter(slot Slot, ft Formatter) {
	if ft == nil {
		panic("nil formatter")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.formatter = ft
	}
}

func (this *logger) SetLinkWriter(slot Slot, wt Writer) {
	if wt == nil {
		panic("nil writer")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.writer = wt
	}
}

func (this *logger) SetLinkLevel(slot Slot, level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.level = level
	}
}

func (this *logger) SetLinkFilter(slot Slot, filter Filter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.filter = filter
	}
}

func (this *logger) link(slot Slot, ft Formatter, wt Writer, opts []interface{}) {
	if ft == nil || wt == nil {
		panic("nil formatter or nil writer")
	}
	link := &link{
		formatter: ft,
		writer:    wt,
		level:     DefaultLevel,
	}
	for _, opt := range opts {
		switch opt := opt.(type) {
		case Level:
			link.level = opt
		case Filter:
			link.filter = opt
		case func(*Record) bool:
			link.filter = opt
		case nil:
			// noop
		default:
			panic("unknown link option type")
		}
	}
	this.slots[slot] = link
}
