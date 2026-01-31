class BalanceSheet {
  constructor() {
    this.balances = new Map();
  }

  addBalance(user, amount) {
    if (this.balances.has(user)) {
      this.balances.set(user, this.balances.get(user) + amount);
    } else {
      this.balances.set(user, amount);
    }
  }

  getBalance(user) {
    return this.balances.get(user) || 0;
  }

  getBalances() {
    return this.balances;
  }
}

module.exports = { BalanceSheet };
