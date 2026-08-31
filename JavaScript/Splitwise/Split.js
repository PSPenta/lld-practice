class Split {
  constructor(user, amount = 0, percentage = 0) {
    if (amount > 0 && percentage > 0) {
      throw new Error('Split cannot have both amount and percentage');
    }

    this.user = user;
    this.amount = amount;
    this.percentage = percentage;
  }
}

module.exports = { Split };
