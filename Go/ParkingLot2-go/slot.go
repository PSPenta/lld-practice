package main

import "fmt"

type Slot2 struct {
	ID          int
	Type        string
	Vehicle     *Vehicle
	IsAvailable bool
}

func NewSlot2(id int, slotType string) (*Slot2, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid slot!")
	}
	return &Slot2{ID: id, Type: slotType, IsAvailable: true}, nil
}

func (s *Slot2) Park(vehicle *Vehicle) bool {
	if !s.IsAvailable {
		fmt.Println("Slot is not empty!")
		return false
	}
	s.Vehicle = vehicle
	s.IsAvailable = false
	return true
}

func (s *Slot2) Unpark() {
	s.Vehicle = nil
	s.IsAvailable = true
}
