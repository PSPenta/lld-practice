const { Balance } = require('./Balance');

class BalanceSheet {
  constructor() {
    this.balances = [];
  }

  addDebt(debtorId, creditorId, amount) {
    if (debtorId === creditorId) return;
    if (!Number.isInteger(amount) || amount <= 0) {
      throw new Error('amount must be a positive integer (paise)');
    }

    const existingBalance = this.balances.find(balance => balance.debtorId === debtorId && balance.creditorId === creditorId);

    const opposite = this.balances.find(balance => balance.debtorId === creditorId && balance.creditorId === debtorId);
    if (opposite && opposite.amount > 0) {
      if (opposite.amount > amount) {
        opposite.amount -= amount;
        return;
      }

      // Remove the opposite balance as it is fully paid
      this.balances.splice(this.balances.indexOf(opposite), 1);
      amount -= opposite.amount;
    }

    if (existingBalance) {
      existingBalance.amount += amount;
      return;
    }

    if (amount <= 0) {
      return;
    }

    this.balances.push(new Balance(debtorId, creditorId, amount));
  }

  getPairwiseBalances() {
    return this.balances;
  }

  getBalance(debtorId, creditorId) {
    return this.balances.find(balance => balance.debtorId === debtorId && balance.creditorId === creditorId);
  }
}

module.exports = { BalanceSheet };
