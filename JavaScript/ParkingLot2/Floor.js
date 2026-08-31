const { Slot } = require('./Slot');

class Floor {
  floor = 0;
  slots = {
    small: [],
    medium: [],
    large: []
  };

  constructor(floor, small, medium, large) {
    if ([floor, small, medium, large].some(num => !num || typeof num != 'number' || num <= 0)) {
      throw new Error('Invalid floor details!');
    }

    this.floor = floor;
    let id = 1;

    for (let i = 0; i < small; i++) {
      this.slots.small.push(new Slot(id++, 'small'));
    }

    for (let i = 0; i < medium; i++) {
      this.slots.medium.push(new Slot(id++, 'medium'));
    }

    for (let i = 0; i < large; i++) {
      this.slots.large.push(new Slot(id++, 'large'));
    }
  }

  findAvailableSlot(type) {
    return this.slots[type].find(slot => slot.isAvailable);
  }

  getSlotById(slotId) {
    if (!slotId || typeof slotId !== 'number') {
      throw new Error('Invalid slot ID!');
    }

    return Object.values(this.slots).flat().find(slot => slot.id == slotId);
  }
}

module.exports = { Floor };
