package main

import (
	"fmt"
	"time"
)

type Ticket struct {
	Floor *Floor
	Slot  *Slot2
	Time  int64
}

func NewTicket(floor *Floor, slot *Slot2) (*Ticket, error) {
	if floor == nil || slot == nil {
		return nil, fmt.Errorf("invalid ticket details!")
	}
	return &Ticket{Floor: floor, Slot: slot, Time: time.Now().UnixMilli()}, nil
}

func (t *Ticket) GetInfo() (floorID, slotID int, ticketTime int64) {
	return t.Floor.FloorNum, t.Slot.ID, t.Time
}
