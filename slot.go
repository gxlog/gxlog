package gxlog

import (
	"fmt"
)

// The Slot defines the slot type of Logger.
type Slot int

// All available slots of Logger here.
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

// MaxSlot is the total count of available slots of Logger.
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

// Link sets the formatter and writer to the slot of Logger. The opts is used
// to specify the slot level and/or the slot filter. An opt must be a value of
// type Level, Filter or func(*Record)bool (the underlying type of the Filter).
//
// If the level of the slot is not specified, it is Trace by default.
// If the filter of the slot is not specified, it is nil by default.
func (this *Logger) Link(slot Slot, formatter Formatter, writer Writer, opts ...interface{}) {
	lnk := link{
		Formatter: formatter,
		Writer:    writer,
		Level:     Trace,
	}
	for _, opt := range opts {
		switch opt := opt.(type) {
		case Level:
			lnk.Level = opt
		case Filter:
			lnk.Filter = opt
		case func(*Record) bool:
			lnk.Filter = opt
		case nil:
			// noop
		default:
			panic(fmt.Sprintf("gxlog.Link: unknown link option type: %T", opt))
		}
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot] = lnk
}

// Unlink sets the formatter, writer and filter of the slot to nil and
// the level of the slot to Off.
func (this *Logger) Unlink(slot Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot] = gNullLink
}

// UnlinkAll sets the formatter, writer and filter of all slots to nil and
// the level of all slots to Off.
func (this *Logger) UnlinkAll() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for i := range this.slots {
		this.slots[i] = gNullLink
	}
}

// CopySlot copies the formatter, writer, level and filter of slot src
// to slot dst.
func (this *Logger) CopySlot(dst, src Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[dst] = this.slots[src]
}

// MoveSlot copies the formatter, writer, level and filter of slot from
// to slot to, and then unlinks the slot from.
func (this *Logger) MoveSlot(to, from Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[to] = this.slots[from]
	this.slots[from] = gNullLink
}

// SwapSlot swaps the formatter, writer, level and filter of the slots.
func (this *Logger) SwapSlot(left, right Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[left], this.slots[right] = this.slots[right], this.slots[left]
}

// SlotFormatter returns the formatter of the slot.
func (this *Logger) SlotFormatter(slot Slot) Formatter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot].Formatter
}

// SetSlotFormatter sets the formatter of the slot.
func (this *Logger) SetSlotFormatter(slot Slot, formatter Formatter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot].Formatter = formatter
}

// SlotWriter returns the writer of the slot.
func (this *Logger) SlotWriter(slot Slot) Writer {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot].Writer
}

// SetSlotWriter sets the writer of the slot.
func (this *Logger) SetSlotWriter(slot Slot, writer Writer) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot].Writer = writer
}

// SlotLevel returns the level of the slot.
func (this *Logger) SlotLevel(slot Slot) Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot].Level
}

// SetSlotLevel sets the level of the slot.
func (this *Logger) SetSlotLevel(slot Slot, level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot].Level = level
}

// SlotFilter returns the filter of the slot.
func (this *Logger) SlotFilter(slot Slot) Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[slot].Filter
}

// SetSlotFilter sets the filter of the slot.
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
