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

func (this *logger) Link(slot Slot, ft Formatter, wt Writer, level Level, filter Filter) {
	this.lock.Lock()

	this.slots[slot] = &link{
		formatter: ft,
		writer:    wt,
		level:     level,
		filter:    filter,
	}
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) Unlink(slot Slot) {
	this.lock.Lock()

	this.slots[slot] = nil
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) UnlinkAll() {
	this.lock.Lock()

	for i := range this.slots {
		this.slots[i] = nil
	}
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) CopyLink(dst, src Slot) {
	this.lock.Lock()

	this.slots[dst] = this.slots[src]
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) MoveLink(to, from Slot) {
	this.lock.Lock()

	this.slots[to] = this.slots[from]
	this.slots[from] = nil
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) SwapLink(left, right Slot) {
	this.lock.Lock()

	this.slots[left], this.slots[right] = this.slots[right], this.slots[left]
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) HasLink(slot Slot) (ok bool) {
	this.lock.Lock()
	ok = (this.slots[slot] != nil)
	this.lock.Unlock()

	return ok
}

func (this *logger) GetLink(slot Slot) (ft Formatter, wt Writer, ok bool) {
	this.lock.Lock()

	lnk := this.slots[slot]
	if lnk != nil {
		ft, wt, ok = lnk.formatter, lnk.writer, true
	} else {
		ft, wt, ok = nil, nil, false
	}

	this.lock.Unlock()

	return ft, wt, ok
}

func (this *logger) MustGetLink(slot Slot) (ft Formatter, wt Writer) {
	this.lock.Lock()
	defer this.lock.Unlock()

	lnk := this.slots[slot]
	return lnk.formatter, lnk.writer
}

func (this *logger) GetLinkLevel(slot Slot) (level Level) {
	this.lock.Lock()

	lnk := this.slots[slot]
	if lnk != nil {
		level = lnk.level
	} else {
		level = LevelOff
	}

	this.lock.Unlock()

	return level
}

func (this *logger) SetLinkLevel(slot Slot, level Level) {
	this.lock.Lock()

	lnk := this.slots[slot]
	if lnk != nil && lnk.level != level {
		lnk.level = level
		this.updateCompactSlots()
	}

	this.lock.Unlock()
}

func (this *logger) GetLinkFilter(slot Slot) (filter Filter) {
	this.lock.Lock()

	lnk := this.slots[slot]
	if lnk != nil {
		filter = lnk.filter
	}

	this.lock.Unlock()

	return filter
}

func (this *logger) SetLinkFilter(slot Slot, filter Filter) {
	this.lock.Lock()

	lnk := this.slots[slot]
	if lnk != nil {
		lnk.filter = filter
	}

	this.lock.Unlock()
}

func (this *logger) updateCompactSlots() {
	this.compactSlots = this.compactSlots[:0]
	for i := range this.slots {
		lnk := this.slots[i]
		if lnk != nil && lnk.level != LevelOff {
			this.compactSlots = append(this.compactSlots, lnk)
		}
	}
}
