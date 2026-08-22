package main

import "fmt"

type Floor struct {
	FloorNum int
	Slots    map[string][]*Slot2
}

func NewFloor(floorNum, small, medium, large int) (*Floor, error) {
	if floorNum <= 0 || small <= 0 || medium <= 0 || large <= 0 {
		return nil, fmt.Errorf("invalid floor details!")
	}
	f := &Floor{
		FloorNum: floorNum,
		Slots: map[string][]*Slot2{
			"small":  {},
			"medium": {},
			"large":  {},
		},
	}
	id := 1
	for i := 0; i < small; i++ {
		f.Slots["small"] = append(f.Slots["small"], mustSlot2(NewSlot2(id, "small")))
		id++
	}
	for i := 0; i < medium; i++ {
		f.Slots["medium"] = append(f.Slots["medium"], mustSlot2(NewSlot2(id, "medium")))
		id++
	}
	for i := 0; i < large; i++ {
		f.Slots["large"] = append(f.Slots["large"], mustSlot2(NewSlot2(id, "large")))
		id++
	}
	return f, nil
}

func (f *Floor) FindAvailableSlot(slotType string) *Slot2 {
	for _, slot := range f.Slots[slotType] {
		if slot.IsAvailable {
			return slot
		}
	}
	return nil
}

func (f *Floor) GetSlotByID(slotID int) *Slot2 {
	for _, slots := range f.Slots {
		for _, slot := range slots {
			if slot.ID == slotID {
				return slot
			}
		}
	}
	return nil
}

func mustSlot2(s *Slot2, err error) *Slot2 {
	if err != nil {
		panic(err)
	}
	return s
}
