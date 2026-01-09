const { ParkingLevel, ParkingLot } = require("./Parkinglot");
const { Slot, slotTypes } = require("./Slot");
const { Car, Bike, Truck } = require("./Vehicle");

const car1 = new Car("car_1234");
const bike1 = new Bike("bike_1234");
const truck1 = new Truck("bike_1234");

const slots = [
  new Slot(1, slotTypes.SMALL),
  new Slot(2, slotTypes.MEDIUM),
  new Slot(3, slotTypes.LARGE),
  new Slot(4, slotTypes.SMALL),
  new Slot(5, slotTypes.MEDIUM),
]

const parkingLevels = [
  new ParkingLevel(slots)
]


const parkingLot = new ParkingLot(parkingLevels);
parkingLot.showAvailability()
parkingLot.park(car1)
parkingLot.showAvailability()



