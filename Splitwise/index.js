const { User } = require('./User');
const { Split } = require('./Split');
const { ExpenseFactory } = require('./Expense');
const { BalanceSheet } = require('./BalanceSheet');

const user1 = new User('XfH4I@example.com');
const user2 = new User('4fNl4@example.com');
const user3 = new User('7Eo5o@example.com');

const balanceSheet = new BalanceSheet();

const split1 = new Split(user1.email);
const split2 = new Split(user2.email);
const split3 = new Split(user3.email);
const expense = ExpenseFactory.createExpense('Equal', user1.email, 300, [
  split1,
  split2,
  split3,
]);
expense.validate();
expense.apply(balanceSheet);
console.log(balanceSheet.getBalances());

const split4 = new Split(user1.email, null, 10);
const split5 = new Split(user2.email, null, 20);
const split6 = new Split(user3.email, null, 70);
const expense2 = ExpenseFactory.createExpense('Percentage', user2.email, 150, [
  split4,
  split5,
  split6,
]);
expense2.validate();
expense2.apply(balanceSheet);
console.log(balanceSheet.getBalances());

const split7 = new Split(user1.email, 225);
const split8 = new Split(user2.email, 390);
const split9 = new Split(user3.email, 680);
const expense3 = ExpenseFactory.createExpense('Exact', user3.email, 1295, [
  split7,
  split8,
  split9,
]);
expense3.validate();
expense3.apply(balanceSheet);
console.log(balanceSheet.getBalances());

console.log(expense);
console.log(expense2);
console.log(expense3);
console.log(balanceSheet.getBalances());
