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

class ExactSplit extends Split {
  constructor(user, amount) {
    super(user, amount);
  }
}

class PercentageSplit extends Split {
  constructor(user, percentage) {
    super(user, 0, percentage);
  }
}

class EqualSplit extends Split {
  constructor(user) {
    super(user);
  }
}

module.exports = { Split, ExactSplit, PercentageSplit, EqualSplit };
