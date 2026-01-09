const { Car, Bike, Truck } = require("./Vehicle");

const slotTypes = {
  SMALL: 'small',
  MEDIUM: 'medium',
  LARGE: 'large',
}

class Slot {
  constructor(slotId, type) {
    this.slotId = slotId;
    if(!Object.values(slotTypes).includes(type)) {
      throw new Error('invalid slot type provided')
    }
    this.type = type;
    this.isAvailable = true;
    this.vehicle = null;
  }

  canFit(vehicle) {
    if(vehicle instanceof Bike) true;

    if(vehicle instanceof Car) {
      return this.type === slotTypes.MEDIUM || this.type === slotTypes.LARGE
    }

    if(vehicle instanceof Truck) return this.type === slotTypes.LARGE

    return false;
  }

  occupy(vehicle) {
    if(!this.isAvailable) {
      throw new Error('slot is unavailable');
    }

    this.isAvailable = false;
    this.vehicle = vehicle;
  }

  release() {
    this.isAvailable = true;
    this.vehicle = null;
  }

}

module.exports = {Slot, slotTypes};
