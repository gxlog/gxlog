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

func (this *logger) Link(slot Slot, ft Formatter, wt Writer, opts ...interface{}) {
	if ft == nil || wt == nil {
		panic("nil formatter or nil writer")
	}

	link := &link{
		formatter: ft,
		writer:    wt,
		level:     LevelTrace,
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

	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot] = link
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

func (this *logger) CopyLink(dst, src Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[dst] = this.slots[src]
}

func (this *logger) MoveLink(to, from Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[to] = this.slots[from]
	this.slots[from] = nil
}

func (this *logger) SwapLink(left, right Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[left], this.slots[right] = this.slots[right], this.slots[left]
}

func (this *logger) Busy(slot Slot) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot] != nil
}

func (this *logger) BusySlots() []Slot {
	this.lock.Lock()
	defer this.lock.Unlock()

	var slots []Slot
	for slot, link := range this.slots {
		if link != nil {
			slots = append(slots, Slot(slot))
		}
	}
	return slots
}

func (this *logger) FreeSlots() []Slot {
	this.lock.Lock()
	defer this.lock.Unlock()

	var slots []Slot
	for slot, link := range this.slots {
		if link == nil {
			slots = append(slots, Slot(slot))
		}
	}
	return slots
}

func (this *logger) LinkFormatter(slot Slot) Formatter {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return nil
	}
	return link.formatter
}

func (this *logger) SetLinkFormatter(slot Slot, ft Formatter) bool {
	if ft == nil {
		panic("nil formatter")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.formatter = ft
		return true
	}
	return false
}

func (this *logger) LinkWriter(slot Slot) Writer {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return nil
	}
	return link.writer
}

func (this *logger) SetLinkWriter(slot Slot, wt Writer) bool {
	if wt == nil {
		panic("nil writer")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.writer = wt
		return true
	}
	return false
}

func (this *logger) LinkLevel(slot Slot) Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return LevelOff
	}
	return link.level
}

func (this *logger) SetLinkLevel(slot Slot, level Level) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.level = level
		return true
	}
	return false
}

func (this *logger) LinkFilter(slot Slot) Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return nil
	}
	return link.filter
}

func (this *logger) SetLinkFilter(slot Slot, filter Filter) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.filter = filter
		return true
	}
	return false
}
