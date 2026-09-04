class User {
  id = 0;
  name = '';
  email = '';

  constructor(id, name, email) {
    if (id <= 0) {
      throw new Error('User ID must be a positive integer');
    }

    if (!name || name.length === 0) {
      throw new Error('User name must be a non-empty string');
    }

    if (!email || !email.includes('@')) {
      throw new Error('User email must be a valid email address');
    }

    this.id = id;
    this.name = name;
    this.email = email;

    return this;
  }
}

module.exports = { User };
