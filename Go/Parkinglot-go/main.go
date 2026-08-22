package main

func main() {
	car1 := NewCar("car_1234")

	slots := []*Slot{
		mustSlot(NewSlot(1, SlotSmall)),
		mustSlot(NewSlot(2, SlotMedium)),
		mustSlot(NewSlot(3, SlotLarge)),
		mustSlot(NewSlot(4, SlotSmall)),
		mustSlot(NewSlot(5, SlotMedium)),
	}

	level, _ := NewParkingLevel(slots)
	parkingLot, _ := NewParkingLot([]*ParkingLevel{level})
	parkingLot.ShowAvailability()
	parkingLot.Park(car1)
	parkingLot.ShowAvailability()
}

func mustSlot(s *Slot, err error) *Slot {
	if err != nil {
		panic(err)
	}
	return s
}
