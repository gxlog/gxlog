package gxlog

import (
	"fmt"
)

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

func (this *Logger) Link(slot Slot, formatter Formatter, writer Writer, opts ...interface{}) {
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
			panic(fmt.Sprintf("gxlog.Link: unknown link option type: %T", opt))
		}
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot] = link
}

func (this *Logger) Unlink(slot Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot] = nil
}

func (this *Logger) UnlinkAll() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for i := range this.slots {
		this.slots[i] = nil
	}
}

func (this *Logger) CopySlot(dst, src Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[src]
	if link == nil {
		this.slots[dst] = nil
	} else {
		copyLink := *link
		this.slots[dst] = &copyLink
	}
}

func (this *Logger) MoveSlot(to, from Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[to] = this.slots[from]
	this.slots[from] = nil
}

func (this *Logger) SwapSlot(left, right Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[left], this.slots[right] = this.slots[right], this.slots[left]
}

func (this *Logger) Busy(slot Slot) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot] != nil
}

func (this *Logger) BusySlots() []Slot {
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

func (this *Logger) FreeSlots() []Slot {
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

func (this *Logger) SlotFormatter(slot Slot) Formatter {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return nil
	}
	return link.Formatter
}

func (this *Logger) SetSlotFormatter(slot Slot, formatter Formatter) bool {
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

func (this *Logger) SlotWriter(slot Slot) Writer {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return nil
	}
	return link.Writer
}

func (this *Logger) SetSlotWriter(slot Slot, writer Writer) bool {
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

func (this *Logger) SlotLevel(slot Slot) Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return LevelOff
	}
	return link.Level
}

func (this *Logger) SetSlotLevel(slot Slot, level Level) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.Level = level
		return true
	}
	return false
}

func (this *Logger) SlotFilter(slot Slot) Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return nil
	}
	return link.Filter
}

func (this *Logger) SetSlotFilter(slot Slot, filter Filter) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link != nil {
		link.Filter = filter
		return true
	}
	return false
}
