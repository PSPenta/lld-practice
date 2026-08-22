package main

import "fmt"

func parkAndPrint(lot *ParkingLot2, v *Vehicle) {
	result, err := lot.ParkVehicle(v)
	if err != nil {
		fmt.Println("Park error:", err)
		return
	}
	fmt.Println("Parked:", result)
}

func unparkAndPrint(lot *ParkingLot2, number string) {
	result, err := lot.UnparkVehicle(number)
	if err != nil {
		fmt.Println("Unpark error:", err)
		return
	}
	fmt.Println("Unparked:", result)
}

func main() {
	floor1, _ := NewFloor(1, 1, 5, 2)
	parkingLot, _ := NewParkingLot2([]*Floor{floor1})
	fmt.Println("Available Slots:", parkingLot.ShowAvailableSlots())

	fmt.Println("Add new floor")
	parkingLot.AddFloor(2, 1, 5, 2)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())

	truck1, _ := CreateVehicle("Truck", "MH12AA2012")
	car1, _ := CreateVehicle("Car", "MH11AA2021")

	parkAndPrint(parkingLot, truck1)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())

	parkAndPrint(parkingLot, car1)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())

	bike1, _ := CreateVehicle("Bike", "MH11AA2022")
	parkAndPrint(parkingLot, bike1)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())

	unparkAndPrint(parkingLot, car1.Number)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())

	bike2, _ := CreateVehicle("Bike", "MH11AA2023")
	parkAndPrint(parkingLot, bike2)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())

	bike3, _ := CreateVehicle("Bike", "MH11AA2024")
	parkAndPrint(parkingLot, bike3)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())

	truck2, _ := CreateVehicle("Truck", "MH12AA2013")
	parkAndPrint(parkingLot, truck2)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())

	truck3, _ := CreateVehicle("Truck", "MH12AA2015")
	parkAndPrint(parkingLot, truck3)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())

	truck4, _ := CreateVehicle("Truck", "MH12AA2016")
	parkAndPrint(parkingLot, truck4)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())

	truck5, _ := CreateVehicle("Truck", "MH12AA2017")
	parkAndPrint(parkingLot, truck5)
	fmt.Println("Available slots", parkingLot.ShowAvailableSlots())
}
