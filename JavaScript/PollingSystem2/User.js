class User {
  id = 0;
  email = '';

  constructor(id, email) {
    if (!id || !email) {
      throw new Error('Invalid user parameters!');
    }

    this.id = id;
    this.email = email;
  }
}

module.exports = { User };
