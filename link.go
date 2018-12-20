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
	Formatter Formatter
	Writer    Writer
	Level     Level
	Filter    Filter
}

func (this *logger) Link(slot Slot, formatter Formatter, writer Writer, opts ...interface{}) {
	if formatter == nil {
		panic("gxlog.Link: nil formatter")
	}
	if writer == nil {
		panic("gxlog.Link: nil writer")
	}

	link := &link{
		Formatter: formatter,
		Writer:    writer,
		Level:     LevelTrace,
	}
	for _, opt := range opts {
		switch opt := opt.(type) {
		case Level:
			link.Level = opt
		case Filter:
			link.Filter = opt
		case func(*Record) bool:
			link.Filter = opt
		case nil:
			// noop
		default:
			panic("gxlog.Link: unknown link option type")
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

func (this *logger) CopySlot(dst, src Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[dst] = this.slots[src]
}

func (this *logger) MoveSlot(to, from Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[to] = this.slots[from]
	this.slots[from] = nil
}

func (this *logger) SwapSlot(left, right Slot) {
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

func (this *logger) SlotFormatter(slot Slot) Formatter {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return nil
	}
	return link.Formatter
}

func (this *logger) SetSlotFormatter(slot Slot, formatter Formatter) bool {
	if formatter == nil {
		panic("gxlog.SetSlotFormatter: nil formatter")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.Formatter = formatter
		return true
	}
	return false
}

func (this *logger) SlotWriter(slot Slot) Writer {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return nil
	}
	return link.Writer
}

func (this *logger) SetSlotWriter(slot Slot, writer Writer) bool {
	if writer == nil {
		panic("gxlog.SetSlotWriter: nil writer")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.Writer = writer
		return true
	}
	return false
}

func (this *logger) SlotLevel(slot Slot) Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return LevelOff
	}
	return link.Level
}

func (this *logger) SetSlotLevel(slot Slot, level Level) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.Level = level
		return true
	}
	return false
}

func (this *logger) SlotFilter(slot Slot) Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return nil
	}
	return link.Filter
}

func (this *logger) SetSlotFilter(slot Slot, filter Filter) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.Filter = filter
		return true
	}
	return false
}
