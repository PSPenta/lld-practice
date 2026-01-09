const { ParkingLot } = require('./ParkingLot');
const { Floor } = require('./Floor');
const { VehicleFactory } = require('./Vehicle');

const parkingLot = new ParkingLot([new Floor(1, 1, 5, 2)]);
console.log('Available Slots:', parkingLot.showAvailableSlots());

console.log('Add new floor');
parkingLot.addFloor(2, 1, 5, 2);
console.log('Available slots', parkingLot.showAvailableSlots());

let truck1 = VehicleFactory.create('Truck', 'MH12AA2012');
let car1 = VehicleFactory.create('Car', 'MH11AA2021');

console.log('Park a truck', parkingLot.parkVehicle(truck1));
console.log('Available slots', parkingLot.showAvailableSlots());

console.log('Park a car', parkingLot.parkVehicle(car1));
console.log('Available slots', parkingLot.showAvailableSlots());

let bike1 = VehicleFactory.create('Bike', 'MH11AA2022');
console.log('Park a bike', parkingLot.parkVehicle(bike1));
console.log('Available slots', parkingLot.showAvailableSlots());

console.log('Unpark a car', parkingLot.unparkVehicle(car1.number));
console.log('Available slots', parkingLot.showAvailableSlots());

let bike2 = VehicleFactory.create('Bike', 'MH11AA2023');
console.log('Park a bike', parkingLot.parkVehicle(bike2));
console.log('Available slots', parkingLot.showAvailableSlots());

let bike3 = VehicleFactory.create('Bike', 'MH11AA2024');
console.log('Park a bike', parkingLot.parkVehicle(bike3));
console.log('Available slots', parkingLot.showAvailableSlots());

let truck2 = VehicleFactory.create('Truck', 'MH12AA2013');
console.log('Park a truck', parkingLot.parkVehicle(truck2));
console.log('Available slots', parkingLot.showAvailableSlots());

let truck3 = VehicleFactory.create('Truck', 'MH12AA2015');
console.log('Park a truck', parkingLot.parkVehicle(truck3));
console.log('Available slots', parkingLot.showAvailableSlots());

let truck4 = VehicleFactory.create('Truck', 'MH12AA2016');
console.log('Park a truck', parkingLot.parkVehicle(truck4));
console.log('Available slots', parkingLot.showAvailableSlots());

let truck5 = VehicleFactory.create('Truck', 'MH12AA2017');
console.log('Park a truck', parkingLot.parkVehicle(truck5));
console.log('Available slots', parkingLot.showAvailableSlots());
