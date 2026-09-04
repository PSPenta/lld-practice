class Balance {
  constructor(debtorId, creditorId, amount) {
    this.debtorId = debtorId;
    this.creditorId = creditorId;
    this.amount = amount;

    return this;
  }
}

module.exports = { Balance };
