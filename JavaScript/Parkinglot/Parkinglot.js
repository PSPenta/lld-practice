class ParkingLevel {
  constructor(slots) {
    if(!slots.length) {
      throw new Error('parking lot slots is missing')
    }
    this.slots = slots || [];
  }

  findSlotToPark(vehicle) {
    return this.slots.find(ele => ele.isAvailable && ele.canFit(vehicle));
  }
}


class ParkingLot {
  constructor(levels) {
    if(!levels.length) {
      throw new Error('parking lot level is missing')
    }
    this.levels = levels || [];
  }

  park(vehicle) {
    let slot = null;
    for (const level of this.levels) {
      slot = level.findSlotToPark(vehicle);
      if(slot) {
        break;
      }
    }

    slot.occupy(vehicle);
  }


  unPark(vehicle) {
    let slot = null;
    for (const level of this.levels) {
      slot = level.slots.find(ele => ele.vehicle.regNumber === vehicle.regNumber);
      if(slot) {
        break;
      }
    }

    if(!slot) {
      throw new Error("vehicle or slot not found")
    }
    slot.release()
  }

  showAvailability() {
    this.levels.forEach((level, index) => {
      const freeSlots = level.slots.filter(slot => slot.isAvailable).length
      console.log(`${freeSlots} slots available on level ${index}`)
    });
  }
}

module.exports = {
  ParkingLevel,
  ParkingLot,
}
