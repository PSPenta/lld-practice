class Vehicle {
  constructor(regNumber) {
    if (new.target === Vehicle) {
      throw new Error("Cannot instantiate abstract class Vehicle directly");
    }
    this.regNumber = regNumber;
  }
}

class Car extends Vehicle {
  constructor(regNumber) {
    super(regNumber)
  }
}

class Bike extends Vehicle {
  constructor(regNumber) {
    super(regNumber)
  }
}

class Truck extends Vehicle {
  constructor(reqNumber) {
    super(reqNumber)
  }
}


module.exports = {
  Bike,
  Car,
  Truck,
}


