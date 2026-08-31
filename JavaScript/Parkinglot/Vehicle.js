/**
 * @interface Vehicle (abstract base class)
 *
 * JavaScript has no native `interface` keyword. Subclasses must be used;
 * `new Vehicle()` throws via `new.target` check.
 *
 * @see Car, Bike, Truck
 */
class Vehicle {
  constructor(regNumber) {
    if (new.target === Vehicle) {
      throw new Error("Cannot instantiate abstract class Vehicle directly");
    }
    this.regNumber = regNumber;
  }
}

/** @implements {Vehicle} */
class Car extends Vehicle {
  constructor(regNumber) {
    super(regNumber)
  }
}

/** @implements {Vehicle} */
class Bike extends Vehicle {
  constructor(regNumber) {
    super(regNumber)
  }
}

/** @implements {Vehicle} */
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
