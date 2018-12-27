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

type slotLink struct {
	Formatter Formatter
	Writer    Writer
	Level     Level
	Filter    Filter
}

var nullSlotLink = slotLink{
	Level: Off,
}

// Link sets the formatter and writer to the slot of Logger. The opts is used
// to specify the slot level and/or the slot filter. An opt must be a value of
// type Level, Filter or func(*Record)bool (the underlying type of the Filter).
//
// If the level of the slot is not specified, it is Trace by default.
// If the filter of the slot is not specified, it is nil by default.
func (log *Logger) Link(slot Slot, formatter Formatter, writer Writer, opts ...interface{}) {
	link := slotLink{
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

	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot] = link
}

// Unlink sets the formatter, writer and filter of the slot to nil and
// the level of the slot to Off.
func (log *Logger) Unlink(slot Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot] = nullSlotLink
}

// UnlinkAll sets the formatter, writer and filter of all slots to nil and
// the level of all slots to Off.
func (log *Logger) UnlinkAll() {
	log.lock.Lock()
	defer log.lock.Unlock()

	for i := range log.slots {
		log.slots[i] = nullSlotLink
	}
}

// CopySlot copies the formatter, writer, level and filter of slot src
// to slot dst.
func (log *Logger) CopySlot(dst, src Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[dst] = log.slots[src]
}

// MoveSlot copies the formatter, writer, level and filter of slot from
// to slot to, and then unlinks the slot from.
func (log *Logger) MoveSlot(to, from Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[to] = log.slots[from]
	log.slots[from] = nullSlotLink
}

// SwapSlot swaps the formatter, writer, level and filter of the slots.
func (log *Logger) SwapSlot(left, right Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[left], log.slots[right] = log.slots[right], log.slots[left]
}

// SlotFormatter returns the formatter of the slot.
func (log *Logger) SlotFormatter(slot Slot) Formatter {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Formatter
}

// SetSlotFormatter sets the formatter of the slot.
func (log *Logger) SetSlotFormatter(slot Slot, formatter Formatter) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Formatter = formatter
}

// SlotWriter returns the writer of the slot.
func (log *Logger) SlotWriter(slot Slot) Writer {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Writer
}

// SetSlotWriter sets the writer of the slot.
func (log *Logger) SetSlotWriter(slot Slot, writer Writer) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Writer = writer
}

// SlotLevel returns the level of the slot.
func (log *Logger) SlotLevel(slot Slot) Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Level
}

// SetSlotLevel sets the level of the slot.
func (log *Logger) SetSlotLevel(slot Slot, level Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Level = level
}

// SlotFilter returns the filter of the slot.
func (log *Logger) SlotFilter(slot Slot) Filter {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Filter
}

// SetSlotFilter sets the filter of the slot.
func (log *Logger) SetSlotFilter(slot Slot, filter Filter) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Filter = filter
}

func (log *Logger) initSlots() {
	for i := 0; i < MaxSlot; i++ {
		log.slots = append(log.slots, nullSlotLink)
	}
}
