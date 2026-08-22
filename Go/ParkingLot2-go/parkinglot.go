package main

import (
	"fmt"
	"strings"
	"time"
)

type Capacity struct {
	Small  int
	Medium int
	Large  int
}

type ParkingLot2 struct {
	Floors           []*Floor
	Capacity         Capacity
	Fare             map[string]float64
	VehicleSlotTypes map[string]string
	Tickets          map[string]*Ticket
}

func NewParkingLot2(floors []*Floor) (*ParkingLot2, error) {
	for _, floor := range floors {
		if floor == nil {
			return nil, fmt.Errorf("invalid floor!")
		}
	}
	p := &ParkingLot2{
		Floors: floors,
		Fare: map[string]float64{
			"small":  1,
			"medium": 2,
			"large":  5,
		},
		VehicleSlotTypes: map[string]string{
			"bike": "small", "scooty": "small",
			"auto": "medium", "car": "medium",
			"truck": "large", "bus": "large",
		},
		Tickets: make(map[string]*Ticket),
	}
	for _, floor := range floors {
		p.Capacity.Small += len(floor.Slots["small"])
		p.Capacity.Medium += len(floor.Slots["medium"])
		p.Capacity.Large += len(floor.Slots["large"])
	}
	return p, nil
}

func (p *ParkingLot2) AddFloor(floorNum, small, medium, large int) error {
	floor, err := NewFloor(floorNum, small, medium, large)
	if err != nil {
		return err
	}
	p.Floors = append(p.Floors, floor)
	p.Capacity.Small += small
	p.Capacity.Medium += medium
	p.Capacity.Large += large
	return nil
}

func (p *ParkingLot2) GetSlotTypeForVehicle(vehicleType string) (string, error) {
	if vehicleType == "" {
		return "", fmt.Errorf("vehicle type must not be empty")
	}
	slotType := p.VehicleSlotTypes[strings.ToLower(vehicleType)]
	if slotType == "" {
		return "", fmt.Errorf("unknown vehicle type: %s", vehicleType)
	}
	return slotType, nil
}

func (p *ParkingLot2) ParkVehicle(vehicle *Vehicle) (map[string]any, error) {
	var slot *Slot2
	var floor *Floor

	for _, f := range p.Floors {
		slot = f.FindAvailableSlot(vehicle.Size)
		if slot != nil {
			floor = f
			break
		}
	}

	if slot == nil && vehicle.Size != "large" {
		for _, f := range p.Floors {
			if vehicle.Size == "medium" {
				slot = f.FindAvailableSlot("large")
				if slot != nil {
					floor = f
					break
				}
			} else {
				if slot == nil {
					slot = f.FindAvailableSlot("large")
					floor = f
				}
				if slot == nil || slot.Type == "large" {
					if mediumSlot := f.FindAvailableSlot("medium"); mediumSlot != nil {
						slot = mediumSlot
						floor = f
						break
					}
				}
			}
		}
	}

	if slot != nil {
		ticket, err := NewTicket(floor, slot)
		if err != nil {
			return nil, fmt.Errorf("could not create ticket: %w", err)
		}
		slot.Park(vehicle)
		p.decrementCapacity(slot.Type)
		p.Tickets[vehicle.Number] = ticket
		_, _, ticketTime := ticket.GetInfo()
		return map[string]any{
			"status":  "Parked",
			"slot":    slot.ID,
			"floor":   floor.FloorNum,
			"vehicle": vehicle.Number,
			"time":    ticketTime,
		}, nil
	}
	return nil, fmt.Errorf("parking full for vehicle type %s", vehicle.Type)
}

func (p *ParkingLot2) decrementCapacity(slotType string) {
	switch slotType {
	case "small":
		p.Capacity.Small--
	case "medium":
		p.Capacity.Medium--
	case "large":
		p.Capacity.Large--
	}
}

func (p *ParkingLot2) incrementCapacity(slotType string) {
	switch slotType {
	case "small":
		p.Capacity.Small++
	case "medium":
		p.Capacity.Medium++
	case "large":
		p.Capacity.Large++
	}
}

func (p *ParkingLot2) UnparkVehicle(vehicleNumber string) (map[string]any, error) {
	if vehicleNumber == "" {
		return nil, fmt.Errorf("vehicle number must not be empty")
	}
	ticket, ok := p.Tickets[vehicleNumber]
	if !ok {
		return nil, fmt.Errorf("vehicle %s is not parked", vehicleNumber)
	}
	floorID, slotID, ticketTime := ticket.GetInfo()
	var floor *Floor
	for _, f := range p.Floors {
		if f.FloorNum == floorID {
			floor = f
			break
		}
	}
	if floor == nil {
		return nil, fmt.Errorf("floor %d not found", floorID)
	}
	slot := floor.GetSlotByID(slotID)
	if slot == nil {
		return nil, fmt.Errorf("slot %d not found on floor %d", slotID, floorID)
	}
	vehicle := slot.Vehicle
	slotType, err := p.GetSlotTypeForVehicle(vehicle.Type)
	if err != nil {
		return nil, err
	}
	slot.Unpark()
	p.incrementCapacity(slotType)
	delete(p.Tickets, vehicleNumber)
	charges := float64(time.Now().UnixMilli()-ticketTime) / 60000.0 * p.Fare[slotType]
	return map[string]any{
		"status":  "Unparked",
		"floor":   floor,
		"slot":    slot,
		"vehicle": vehicleNumber,
		"charges": charges,
	}, nil
}

func (p *ParkingLot2) ShowAvailableSlots() Capacity {
	return p.Capacity
}

func (p *ParkingLot2) ShowParkedVehicles() map[string]*Ticket {
	return p.Tickets
}
