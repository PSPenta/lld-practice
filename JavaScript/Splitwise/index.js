const { SplitwiseService } = require('./SplitwiseService');
const { EqualSplit, PercentageSplit, ExactSplit } = require('./Split');

const service = new SplitwiseService();

const user1 = service.addUser('John Doe', 'john.doe@example.com');
try {
  const user2 = service.addUser('John Doe', 'john.doe@example.com');
} catch (error) {
  console.log('Error adding user: ', error.message);
}

const user2 = service.addUser('Jane Doe', 'jane.doe@example.com');
const user3 = service.addUser('Jim Doe', 'jim.doe@example.com');

const expense = service.addExpense({
  type: 'Equal',
  paidBy: user1.id,
  amount: 300,
  splits: [
    new EqualSplit(user1.id),
    new EqualSplit(user2.id),
    new EqualSplit(user3.id),
  ],
});
console.log(service.getPairwiseBalances());

const expense2 = service.addExpense({
  type: 'Percentage',
  paidBy: user2.id,
  amount: 150.50,
  splits: [
    new PercentageSplit(user1.id, 10),
    new PercentageSplit(user2.id, 20),
    new PercentageSplit(user3.id, 70),
  ],
});
console.log(service.getPairwiseBalances());

service.settleUp(user2.id, user1.id);
console.log(service.getPairwiseBalances());

service.settleUp(user3.id, user2.id);
console.log(service.getPairwiseBalances());

service.settleUp(user3.id, user1.id, 99);
console.log(service.getPairwiseBalances());

try {
  service.addExpense({
    type: 'Exact',
    paidBy: user3.id,
    amount: 100,
    splits: [
      new ExactSplit(user1.id, 30),
      new ExactSplit(user2.id, 20),
      new ExactSplit(user3.id, 51),
    ],
  });
} catch (error) {
  console.log('Error adding expense: ', error.message);
}
console.log(service.getPairwiseBalances());

console.log('\nFinal settlements');
service.settleUp(user3.id, user1.id);
console.log(service.getPairwiseBalances());
