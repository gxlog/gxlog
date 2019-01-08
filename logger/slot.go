package logger

import (
	"fmt"

	"github.com/gxlog/gxlog/iface"
)

// The Slot defines the slot type of Logger.
type Slot int

// All available slots here.
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

// MaxSlot is the total count of slots.
const MaxSlot = 8

type slotLink struct {
	Formatter iface.Formatter
	Writer    iface.Writer
	Level     iface.Level
	Filter    Filter
}

var nullSlotLink = slotLink{
	Level: iface.Off,
}

// Link sets the formatter and writer to the slot. The opts is used to specify
// the slot Level and/or the slot Filter. An opt MUST be a value of type Level,
// Filter or func(*Record)bool (the underlying type of Filter).
// If the Level of the slot is not specified, Trace is used.
func (log *Logger) Link(slot Slot, formatter iface.Formatter,
	writer iface.Writer, opts ...interface{}) {
	link := slotLink{
		Formatter: formatter,
		Writer:    writer,
		Level:     iface.Trace,
	}
	for _, opt := range opts {
		switch opt := opt.(type) {
		case iface.Level:
			link.Level = opt
		case Filter:
			link.Filter = opt
		case func(*iface.Record) bool:
			link.Filter = opt
		case nil:
			// noop
		default:
			panic(fmt.Sprintf("logger.Link: unknown link option type: %T", opt))
		}
	}

	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot] = link
}

// Unlink sets the Formatter, Writer and Filter of the slot to nil and
// the Level of the slot to Off.
func (log *Logger) Unlink(slot Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot] = nullSlotLink
}

// UnlinkAll sets the Formatter, Writer and Filter of all slots to nil and
// the Level of all slots to Off.
func (log *Logger) UnlinkAll() {
	log.lock.Lock()
	defer log.lock.Unlock()

	for i := range log.slots {
		log.slots[i] = nullSlotLink
	}
}

// CopySlot copies the Formatter, Writer, Level and Filter of Slot src
// to Slot dst.
func (log *Logger) CopySlot(dst, src Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[dst] = log.slots[src]
}

// MoveSlot copies the Formatter, Writer, Level and Filter of Slot from
// to Slot to, and then unlinks Slot from.
func (log *Logger) MoveSlot(to, from Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[to] = log.slots[from]
	log.slots[from] = nullSlotLink
}

// SwapSlot swaps the Formatter, Writer, Level and Filter of the slots.
func (log *Logger) SwapSlot(left, right Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[left], log.slots[right] = log.slots[right], log.slots[left]
}

// SlotFormatter returns the Formatter of the slot.
func (log *Logger) SlotFormatter(slot Slot) iface.Formatter {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Formatter
}

// SetSlotFormatter sets the Formatter of the slot.
func (log *Logger) SetSlotFormatter(slot Slot, formatter iface.Formatter) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Formatter = formatter
}

// SlotWriter returns the Writer of the slot.
func (log *Logger) SlotWriter(slot Slot) iface.Writer {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Writer
}

// SetSlotWriter sets the Writer of the slot.
func (log *Logger) SetSlotWriter(slot Slot, writer iface.Writer) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Writer = writer
}

// SlotLevel returns the Level of the slot.
func (log *Logger) SlotLevel(slot Slot) iface.Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Level
}

// SetSlotLevel sets the Level of the slot.
func (log *Logger) SetSlotLevel(slot Slot, level iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Level = level
}

// SlotFilter returns the Filter of the slot.
func (log *Logger) SlotFilter(slot Slot) Filter {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Filter
}

// SetSlotFilter sets the Filter of the slot.
func (log *Logger) SetSlotFilter(slot Slot, filter Filter) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Filter = filter
}

func (log *Logger) initSlots() {
	log.slots = make([]slotLink, MaxSlot)
	for i := range log.slots {
		log.slots[i] = nullSlotLink
	}
}
