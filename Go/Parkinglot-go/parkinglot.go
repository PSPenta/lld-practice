package main

import "fmt"

type ParkingLevel struct {
	Slots []*Slot
}

func NewParkingLevel(slots []*Slot) (*ParkingLevel, error) {
	if len(slots) == 0 {
		return nil, fmt.Errorf("parking lot slots is missing")
	}
	return &ParkingLevel{Slots: slots}, nil
}

func (l *ParkingLevel) FindSlotToPark(vehicle *Vehicle) *Slot {
	for _, slot := range l.Slots {
		if slot.IsAvailable && slot.CanFit(vehicle) {
			return slot
		}
	}
	return nil
}

type ParkingLot struct {
	Levels []*ParkingLevel
}

func NewParkingLot(levels []*ParkingLevel) (*ParkingLot, error) {
	if len(levels) == 0 {
		return nil, fmt.Errorf("parking lot level is missing")
	}
	return &ParkingLot{Levels: levels}, nil
}

func (p *ParkingLot) Park(vehicle *Vehicle) {
	var slot *Slot
	for _, level := range p.Levels {
		slot = level.FindSlotToPark(vehicle)
		if slot != nil {
			break
		}
	}
	slot.Occupy(vehicle)
}

func (p *ParkingLot) UnPark(vehicle *Vehicle) error {
	var slot *Slot
	for _, level := range p.Levels {
		for _, s := range level.Slots {
			if s.Vehicle != nil && s.Vehicle.RegNumber == vehicle.RegNumber {
				slot = s
				break
			}
		}
		if slot != nil {
			break
		}
	}
	if slot == nil {
		return fmt.Errorf("vehicle or slot not found")
	}
	slot.Release()
	return nil
}

func (p *ParkingLot) ShowAvailability() {
	for index, level := range p.Levels {
		free := 0
		for _, slot := range level.Slots {
			if slot.IsAvailable {
				free++
			}
		}
		fmt.Printf("%d slots available on level %d\n", free, index)
	}
}
