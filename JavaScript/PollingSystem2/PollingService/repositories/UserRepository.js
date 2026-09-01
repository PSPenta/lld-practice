const { User } = require('../models/User');

/** In-memory stand-in for a user repository. */
class UserRepository {
  static users = [];
  static nextId = 1;

  static createId() {
    return this.nextId++;
  }

  static add(user) {
    if (!(user instanceof User)) {
      throw new Error('Invalid user!');
    }

    if (this.users.some((u) => u.id === user.id)) {
      throw new Error('User already exists!');
    }

    this.users.push(user);
  }

  static create(email = '') {
    if (!email) {
      throw new Error('Invalid user email!');
    }

    if (this.users.some((u) => u.email === email)) {
      throw new Error('User already exists!');
    }

    const user = new User(this.createId(), email);
    this.add(user);
    return user;
  }

  static getById(id) {
    return this.users.find((u) => u.id === id);
  }
}

module.exports = { UserRepository };
