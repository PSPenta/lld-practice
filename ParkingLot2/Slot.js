class Slot {
  id = 0;
  type = '';
  vehicle = null;
  isAvailable = true;

  constructor(id, type) {
    if (!id || typeof id !== 'number' || id <= 0) {
      throw new Error('Invalid slot!');
    }

    this.id = id;
    this.type = type;
  }

  park(vehicle) {
    if (!this.isAvailable) {
      console.log('Slot is not empty!');
      return false;
    }

    this.vehicle = vehicle;
    this.isAvailable = false;

    return true;
  }

  unpark() {
    this.vehicle = null;
    this.isAvailable = true;
  }
}

module.exports = { Slot };
