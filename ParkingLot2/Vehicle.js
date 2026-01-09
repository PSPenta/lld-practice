class Vehicle {
  type = '';
  number = '';
  size = '';

  constructor(type, number, size) {
    if (!type || !number || !size) {
      throw new Error('Invalid vehicle!');
    }
    this.type = type;
    this.size = size;
    this.number = number;
  }
}

class Bike extends Vehicle {
  constructor(number, size) {
    super('Bike', number, size);
  }
}

class Car extends Vehicle {
  constructor(number, size) {
    super('Car', number, size);
  }
}

class Truck extends Vehicle {
  constructor(number, size) {
    super('Truck', number, size);
  }
}

class VehicleFactory {
  static create(type, number) {
    if (!type || !number) {
      throw new Error('Invalid vehicle inputs!');
    }

    const normalized = type.toLowerCase();

    switch (normalized) {
      case 'bike':  return new Bike(number, 'small');
      case 'car':   return new Car(number, 'medium');
      case 'truck': return new Truck(number, 'large');

      default:
        throw new Error(`Unsupported vehicle type: ${type}`);
    }
  }
}

module.exports = { VehicleFactory };
