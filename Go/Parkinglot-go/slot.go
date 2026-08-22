package main

import "fmt"

type SlotType string

const (
	SlotSmall  SlotType = "small"
	SlotMedium SlotType = "medium"
	SlotLarge  SlotType = "large"
)

type Slot struct {
	SlotID      int
	Type        SlotType
	IsAvailable bool
	Vehicle     *Vehicle
}

func NewSlot(slotID int, slotType SlotType) (*Slot, error) {
	switch slotType {
	case SlotSmall, SlotMedium, SlotLarge:
	default:
		return nil, fmt.Errorf("invalid slot type provided")
	}
	return &Slot{SlotID: slotID, Type: slotType, IsAvailable: true}, nil
}

func (s *Slot) CanFit(vehicle *Vehicle) bool {
	if vehicle.Type == VehicleBike {
		return true
	}
	if vehicle.Type == VehicleCar {
		return s.Type == SlotMedium || s.Type == SlotLarge
	}
	if vehicle.Type == VehicleTruck {
		return s.Type == SlotLarge
	}
	return false
}

func (s *Slot) Occupy(vehicle *Vehicle) error {
	if !s.IsAvailable {
		return fmt.Errorf("slot is unavailable")
	}
	s.IsAvailable = false
	s.Vehicle = vehicle
	return nil
}

func (s *Slot) Release() {
	s.IsAvailable = true
	s.Vehicle = nil
}
