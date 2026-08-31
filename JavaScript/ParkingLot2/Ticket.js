class Ticket {
  floor = 0;
  slot = 0;
  time = new Date();

  constructor(floor, slot) {
    if ([floor, slot].some(ele => !ele)) {
console.log(floor, slot);
      throw new Error('Invalid ticket details!');
    }

    this.floor = floor;
    this.slot = slot;
    this.time = Date.now();
  }

  getInfo() {
    return {
      floorId: this.floor,
      slotId: this.slot,
      time: this.time
    };
  }
}

module.exports = { Ticket };
