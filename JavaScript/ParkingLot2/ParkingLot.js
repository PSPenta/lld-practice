const { Floor } = require('./Floor');
const { Ticket } = require('./Ticket');

class ParkingLot {
  floors = [];

  capacity = {
    small: 0,
    medium: 0,
    large: 0
  };

  fare = {
    small: 1,
    medium: 2,
    large: 5
  };

  vehicleSlotTypes = {
    bike: 'small',
    scooty: 'small',
    auto: 'medium',
    car: 'medium',
    truck: 'large',
    bus: 'large',
  }

  tickets = new Map();

  constructor(floors) {
    if (
      floors.some(floor => !(floor instanceof Floor))
    ) {
      throw new Error('Invalid floor!');
    }

    this.floors = floors;

    for (let i = 0; i < floors.length; i++) {
      this.capacity.small = floors[i].slots.small.length;
      this.capacity.medium = floors[i].slots.medium.length;
      this.capacity.large = floors[i].slots.large.length;
    }
  }

  addFloor(floor, small, medium, large) {
    this.floors.push(new Floor(floor, small, medium, large));

    this.capacity.small += small;
    this.capacity.medium += medium;
    this.capacity.large += large;

    return true;
  }

  getSlotTypeForVehicle(type) {
    if (!type) {
      throw new Error('Invalid vehicle type!');
    }

    return this.vehicleSlotTypes[type.toLowerCase()];
  }

  parkVehicle(vehicle) {
    if (!vehicle) {
      throw new Error('Invalid vehicle!');
    }

    let slot = null, Floor = null;
    for (let floor of this.floors) {
      slot = floor.findAvailableSlot(vehicle.size);
      if (slot) {
        Floor = floor;
        break;
      }
    }

    if (!slot && vehicle.size != 'large') {
      for (let floor of this.floors) {
        if (vehicle.size == 'medium') {
          slot = floor.findAvailableSlot('large');
          if (slot) {
            Floor = floor;
            break;
          }
        } else {
          if (!slot) {
            slot = floor.findAvailableSlot('large');
            Floor = floor;
          }

          if (!slot || slot.type == 'large') {
            let mediumSlot = floor.findAvailableSlot('medium');
            if (mediumSlot) {
              slot = mediumSlot;
              Floor = floor;
              break;
            }
          }
        }
      }
    }

    if (slot) {
      let ticket = new Ticket(Floor.floor, slot.id);
      slot.park(vehicle);
      this.capacity[slot.type]--;
      this.tickets.set(vehicle.number, ticket);

      return {
        status: 'Parked',
        slot: slot.id,
        floor: Floor.floor,
        vehicle: vehicle.number,
        time: ticket.getInfo().time
      }
    }

    throw new Error(`Parking full for vehicle ${vehicle.type}!`);
  }

  unparkVehicle(vehicleNumber) {
    if (!vehicleNumber || typeof vehicleNumber !== 'string') {
      throw new Error('Invalid vehicle!');
    }

    let ticket = this.tickets.get(vehicleNumber);
    if (!ticket || !(ticket instanceof Ticket)) {
      throw new Error('Vehicle is not parked!');
    }

    let { floorId, slotId, time } = ticket.getInfo();
    let floor = this.floors.find(floor => floor.floor == floorId);
    if (!floor) {
      throw new Error('Invalid floor!');
    }

    let slot = floor.getSlotById(slotId);
    if (!slot) {
      throw new Error('Invalid slot!');
    }

    let { vehicle } = slot;
    let slotType = this.getSlotTypeForVehicle(vehicle.type);
    slot.unpark();
    this.capacity[slotType]++;
    this.tickets.delete(vehicleNumber);

    return {
      status: 'Unparked',
      floor: floor,
      slot: slot,
      vehicle: vehicleNumber,
      charges: ((Date.now() - time) / 60000) * this.fare[slotType]
    };
  }

  showAvailableSlots() {
    return this.capacity;
  }

  showParkedVehicles() {
    return this.tickets;
  }
}

module.exports = { ParkingLot };
