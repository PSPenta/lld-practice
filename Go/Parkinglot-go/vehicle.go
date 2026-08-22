package main

type VehicleType int

const (
	VehicleBike VehicleType = iota
	VehicleCar
	VehicleTruck
)

type Vehicle struct {
	Type      VehicleType
	RegNumber string
}

func NewCar(regNumber string) *Vehicle {
	return &Vehicle{Type: VehicleCar, RegNumber: regNumber}
}

func NewBike(regNumber string) *Vehicle {
	return &Vehicle{Type: VehicleBike, RegNumber: regNumber}
}

func NewTruck(regNumber string) *Vehicle {
	return &Vehicle{Type: VehicleTruck, RegNumber: regNumber}
}
