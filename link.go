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

type Link struct {
	Formatter Formatter
	Writer    Writer
	Level     Level
	Filter    Filter
}

func NewLink(ft Formatter, wt Writer) *Link {
	return &Link{
		Formatter: ft,
		Writer:    wt,
	}
}

func (this *Link) WithLevel(level Level) *Link {
	this.Level = level
	return this
}

func (this *Link) WithFilter(filter Filter) *Link {
	this.Filter = filter
	return this
}

func (this *logger) GetLink(slot Slot) *Link {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return nil
	}
	copyLink := *link
	return &copyLink
}

func (this *logger) SetLink(slot Slot, link *Link) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if link == nil {
		this.slots[slot] = nil
		return
	}
	this.setLink(slot, link)
}

func (this *logger) UpdateLink(slot Slot, fn func(*Link)) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	link := this.slots[slot]
	if link == nil {
		return false
	}
	copyLink := *link
	fn(&copyLink)
	this.setLink(slot, &copyLink)
	return true
}

func (this *logger) ResetLink(slot Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[slot] = nil
}

func (this *logger) ClearLinks() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for i := range this.slots {
		this.slots[i] = nil
	}
}

func (this *logger) setLink(slot Slot, link *Link) {
	copyLink := *link
	this.slots[slot] = &copyLink
}
