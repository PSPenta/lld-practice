const { Admin } = require('./Admin');
const { Polls } = require('./Polls');
const { User } = require('./User');

let admin = new Admin(1);

let poll = admin.createPoll('Whats capital of India?', ['Delhi', 'Mumbai']);
Polls.addPoll(poll);

let user = new User(2);
user.submitPoll(poll, 'Mumbai');

let user2 = new User(3);
user2.submitPoll(poll, 'Delhi');

let user3 = new User(3);
user3.submitPoll(poll, 'Delhi');

let user4 = new User(3);
user4.submitPoll(poll, 'Delhi');

admin.showStatistics(poll);

