const { User } = require('./User');
const { Split } = require('./Split');
const { ExpenseFactory } = require('./Expense');

const user1 = new User('XfH4I@example.com');
const user2 = new User('4fNl4@example.com');
const user3 = new User('7Eo5o@example.com');

const split1 = new Split(user1);
const split2 = new Split(user2);
const split3 = new Split(user3);
const expense = ExpenseFactory.createExpense('Equal', user1, 300, [
  split1,
  split2,
  split3,
]);
expense.validate();

const split4 = new Split(user1, null, 10);
const split5 = new Split(user2, null, 20);
const split6 = new Split(user3, null, 70);
const expense2 = ExpenseFactory.createExpense('Percentage', user2, 150, [
  split4,
  split5,
  split6,
]);
expense2.validate();

const split7 = new Split(user1, 225);
const split8 = new Split(user2, 390);
const split9 = new Split(user3, 680);
const expense3 = ExpenseFactory.createExpense('Exact', user3, 1295, [
  split7,
  split8,
  split9,
]);
expense3.validate();

console.log(expense);
console.log(expense2);
console.log(expense3);
