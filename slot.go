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
)

const MaxSlot = 8

type link struct {
	Formatter Formatter
	Writer    Writer
	Level     Level
	Filter    Filter
}

var gNullLink = link{
	Level: Off,
}

func (this *Logger) Link(slot Slot, formatter Formatter, writer Writer, opts ...interface{}) {
	link := link{
		Formatter: formatter,
		Writer:    writer,
		Level:     Trace,
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

	this.slots[slot] = gNullLink
}

func (this *Logger) UnlinkAll() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for i := range this.slots {
		this.slots[i] = gNullLink
	}
}

func (this *Logger) CopySlot(dst, src Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[dst] = this.slots[src]
}

func (this *Logger) MoveSlot(to, from Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[to] = this.slots[from]
	this.slots[from] = gNullLink
}

func (this *Logger) SwapSlot(left, right Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[left], this.slots[right] = this.slots[right], this.slots[left]
}

func (this *Logger) SlotFormatter(slot Slot) Formatter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot].Formatter
}

func (this *Logger) SetSlotFormatter(slot Slot, formatter Formatter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot].Formatter = formatter
}

func (this *Logger) SlotWriter(slot Slot) Writer {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot].Writer
}

func (this *Logger) SetSlotWriter(slot Slot, writer Writer) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot].Writer = writer
}

func (this *Logger) SlotLevel(slot Slot) Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot].Level
}

func (this *Logger) SetSlotLevel(slot Slot, level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot].Level = level
}

func (this *Logger) SlotFilter(slot Slot) Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot].Filter
}

func (this *Logger) SetSlotFilter(slot Slot, filter Filter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot].Filter = filter
}

func (this *Logger) initSlots() {
	for i := 0; i < MaxSlot; i++ {
		this.slots = append(this.slots, gNullLink)
	}
}
