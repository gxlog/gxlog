package gxlog

type link struct {
	formatter Formatter
	writer    Writer
	level     LogLevel
}

func (this *logger) Link(slot LinkSlot, ft Formatter, wt Writer, level LogLevel) {
	this.lock.Lock()

	this.linkSlots[slot] = &link{
		formatter: ft,
		writer:    wt,
		level:     level,
	}
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) Unlink(slot LinkSlot) {
	this.lock.Lock()

	this.linkSlots[slot] = nil
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) UnlinkAll() {
	this.lock.Lock()

	for i := range this.linkSlots {
		this.linkSlots[i] = nil
	}
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) CopyLink(src, dst LinkSlot) {
	this.lock.Lock()

	this.linkSlots[dst] = this.linkSlots[src]
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) MoveLink(from, to LinkSlot) {
	this.lock.Lock()

	this.linkSlots[to] = this.linkSlots[from]
	this.linkSlots[from] = nil
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) SwapLink(left, right LinkSlot) {
	this.lock.Lock()

	this.linkSlots[left], this.linkSlots[right] = this.linkSlots[right], this.linkSlots[left]
	this.updateCompactSlots()

	this.lock.Unlock()
}

func (this *logger) HasLink(slot LinkSlot) (ok bool) {
	this.lock.Lock()
	ok = (this.linkSlots[slot] != nil)
	this.lock.Unlock()

	return ok
}

func (this *logger) GetLink(slot LinkSlot) (ft Formatter, wt Writer, ok bool) {
	this.lock.Lock()

	lnk := this.linkSlots[slot]
	if lnk != nil {
		ft, wt, ok = lnk.formatter, lnk.writer, true
	} else {
		ft, wt, ok = nil, nil, false
	}

	this.lock.Unlock()

	return ft, wt, ok
}

func (this *logger) MustGetLink(slot LinkSlot) (ft Formatter, wt Writer) {
	this.lock.Lock()
	defer this.lock.Unlock()

	lnk := this.linkSlots[slot]
	return lnk.formatter, lnk.writer
}

func (this *logger) GetLinkLevel(slot LinkSlot) (level LogLevel) {
	this.lock.Lock()

	lnk := this.linkSlots[slot]
	if lnk != nil {
		level = lnk.level
	} else {
		level = LevelOff
	}

	this.lock.Unlock()

	return level
}

func (this *logger) SetLinkLevel(slot LinkSlot, level LogLevel) {
	this.lock.Lock()

	lnk := this.linkSlots[slot]
	if lnk != nil && lnk.level != level {
		lnk.level = level
		this.updateCompactSlots()
	}

	this.lock.Unlock()
}

func (this *logger) updateCompactSlots() {
	this.compactSlots = this.compactSlots[:0]
	for i := range this.linkSlots {
		lnk := this.linkSlots[i]
		if lnk != nil && lnk.level != LevelOff {
			this.compactSlots = append(this.compactSlots, lnk)
		}
	}
}
