package main

import "fmt"

type Vehicle struct {
	Type   string
	Number string
	Size   string
}

func NewBike(number, size string) *Vehicle {
	return &Vehicle{Type: "Bike", Number: number, Size: size}
}

func NewCar(number, size string) *Vehicle {
	return &Vehicle{Type: "Car", Number: number, Size: size}
}

func NewTruck(number, size string) *Vehicle {
	return &Vehicle{Type: "Truck", Number: number, Size: size}
}

func CreateVehicle(vehicleType, number string) (*Vehicle, error) {
	if vehicleType == "" || number == "" {
		return nil, fmt.Errorf("invalid vehicle inputs!")
	}
	switch vehicleType {
	case "bike", "Bike":
		return NewBike(number, "small"), nil
	case "car", "Car":
		return NewCar(number, "medium"), nil
	case "truck", "Truck":
		return NewTruck(number, "large"), nil
	default:
		return nil, fmt.Errorf("unsupported vehicle type: %s", vehicleType)
	}
}
