package gxlog

type link struct {
	ft     Formatter
	wt     Writer
	enable bool
}

func (this *Logger) Link(ft Formatter, wt Writer, slot LinkSlot) {
	this.linkSlots[slot] = &link{ft, wt, true}
	this.updateCompactSlots()
}

func (this *Logger) Unlink(slot LinkSlot) {
	this.linkSlots[slot] = nil
	this.updateCompactSlots()
}

func (this *Logger) UnlinkAll() {
	for i := range this.linkSlots {
		this.linkSlots[i] = nil
	}
	this.updateCompactSlots()
}

func (this *Logger) CopyLink(src, dst LinkSlot) {
	this.linkSlots[dst] = this.linkSlots[src]
	this.updateCompactSlots()
}

func (this *Logger) MoveLink(from, to LinkSlot) {
	this.linkSlots[to] = this.linkSlots[from]
	this.linkSlots[from] = nil
	this.updateCompactSlots()
}

func (this *Logger) SwapLink(left, right LinkSlot) {
	this.linkSlots[left], this.linkSlots[right] = this.linkSlots[right], this.linkSlots[left]
	this.updateCompactSlots()
}

func (this *Logger) HasLink(slot LinkSlot) bool {
	return this.linkSlots[slot] != nil
}

func (this *Logger) GetLink(slot LinkSlot) (Formatter, Writer, bool) {
	link := this.linkSlots[slot]
	if link == nil {
		return nil, nil, false
	}
	return link.ft, link.wt, true
}

func (this *Logger) EnableLink(slot LinkSlot) {
	this.setLinkEnable(slot, true)
}

func (this *Logger) DisableLink(slot LinkSlot) {
	this.setLinkEnable(slot, false)
}

func (this *Logger) setLinkEnable(slot LinkSlot, enable bool) {
	link := this.linkSlots[slot]
	if link != nil && link.enable != enable {
		link.enable = enable
		this.updateCompactSlots()
	}
}

func (this *Logger) updateCompactSlots() {
	this.compactSlots = this.compactSlots[:0]
	for i := range this.linkSlots {
		link := this.linkSlots[i]
		if link != nil && link.enable {
			this.compactSlots = append(this.compactSlots, link)
		}
	}
}
