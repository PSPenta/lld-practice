const { User } = require('./User');
const { ExpenseFactory } = require('./Expense');
const { BalanceSheet } = require('./BalanceSheet');
const { toAmount, fromAmount } = require('./money');

class SplitwiseService {
  constructor() {
    this.userId = 0;
    this.expenseId = 0;
    this.users = new Map();
    this.expenses = new Map();
    this.balanceSheet = new BalanceSheet();
  }

  addUser(name, email) {
    if (!name || !email || !email.includes('@') || !email.includes('.')) {
      throw new Error('Invalid user input');
    }

    if (Array.from(this.users.values()).some(user => user.email === email)) {
      throw new Error('User already exists');
    }

    const user = new User(++this.userId, name, email);
    this.users.set(this.userId, user);

    return user;
  }

  getUser(id) {
    const user = this.users.get(id);
    if (!user) {
      throw new Error('User not found');
    }

    return user;
  }

  addExpense(input) {
    const { type, paidBy, amount, splits } = input;
    if (!type || !paidBy || !amount || !splits) {
      throw new Error('Invalid expense input');
    }

    if (!this.users.has(paidBy)) {
      throw new Error('User not found');
    }

    if (splits.some(split => !this.users.has(split.user))) {
      throw new Error('User not found');
    }

    // convert rupees → paise once at the boundary; ledger stays integer
    const amountPaise = toAmount(amount);
    for (const split of splits) {
      if (split.amount) split.amount = toAmount(split.amount);
    }

    const expense = ExpenseFactory.createExpense({
      ...input,
      id: ++this.expenseId,
      amount: amountPaise,
      splits,
    });
    this.expenses.set(this.expenseId, expense);

    expense.validate();
    expense.apply(this.balanceSheet);

    return expense;
  }

  getPairwiseBalances() {
    const pairwiseBalances = this.balanceSheet.getPairwiseBalances();
    const result = [];
    pairwiseBalances.forEach(balance => {
      result.push(
        `User ${this.users.get(balance.debtorId).email} (${balance.debtorId}) owes User ${this.users.get(balance.creditorId).email} (${balance.creditorId}) ${fromAmount(balance.amount)}`,
      );
    });

    if (result.length === 0) {
      return 'All balances are settled!';
    }

    return result;
  }

  settleUp(payerId, payeeId, amount = null) {
    const payer = this.users.get(payerId);
    const payee = this.users.get(payeeId);

    if (!payer || !payee) {
      throw new Error('User not found');
    }

    const balance = this.balanceSheet.getBalance(payerId, payeeId);
    if (!balance || balance.amount <= 0) {
      throw new Error('Balance not found');
    }

    const pay =
      amount == null ? balance.amount : toAmount(amount);
    if (!Number.isInteger(pay) || pay <= 0) {
      throw new Error('invalid settle amount');
    }
    if (pay > balance.amount) {
      throw new Error('Amount is greater than balance');
    }

    this.balanceSheet.addDebt(payeeId, payerId, pay);
    console.log(
      `User ${payer.email} (${payerId}) settled up with User ${payee.email} (${payeeId}) ${fromAmount(pay)}`,
    );
  }
}

module.exports = { SplitwiseService };
