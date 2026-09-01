const { PollService } = require('./PollService');

const pollService = new PollService();

const alice = pollService.createUser('alice@example.com');
const bob = pollService.createUser('bob@example.com');
const jane = pollService.createUser('jane@example.com');
const jim = pollService.createUser('jim@example.com');
const kate = pollService.createUser('kate@example.com');

const alicePoll = pollService.createPoll(
  alice,
  'Whats capital of India?',
  ['Delhi', 'Mumbai'],
  true,
  false,
);

try {
  pollService.assignVoter(alice, alicePoll, alice);
} catch (err) {
  console.log('Expected error:', err.message);
}

try {
  pollService.submitVote(alice, alicePoll, 'Delhi');
} catch (err) {
  console.log('Expected error:', err.message);
}

try {
  pollService.submitVote(bob, alicePoll, 'Mumbai');
} catch (err) {
  console.log('Expected error:', err.message);
}

pollService.assignVoter(alice, alicePoll, bob);
pollService.submitVote(bob, alicePoll, 'Mumbai');

pollService.assignVoter(alice, alicePoll, jane);
pollService.submitVote(jane, alicePoll, 'Delhi');

pollService.assignVoter(alice, alicePoll, jim);
pollService.submitVote(jim, alicePoll, 'Delhi');

alicePoll.isClosed = true;
pollService.updatePoll(alicePoll);

try {
  pollService.assignVoter(alice, alicePoll, kate);
  pollService.submitVote(kate, alicePoll, 'Delhi');
} catch (err) {
  console.log('Expected error:', err.message);
}

try {
  pollService.submitVote(jim, alicePoll, 'Delhi');
} catch (err) {
  console.log('Expected error:', err.message);
}

console.log(
  'Poll statistics for',
  alicePoll.id,
  ':',
  pollService.getStatistics(alice, alicePoll),
);

// Alice can vote on Bob's poll (creator rule is per-poll)
const bobPoll = pollService.createPoll(bob, 'Best language?', ['JavaScript', 'Go'], false, false);
// pollService.assignVoter(bob, bobPoll, alice);
pollService.submitVote(alice, bobPoll, 'JavaScript');
console.log(
  'Poll statistics for',
  bobPoll.id,
  ':',
  pollService.getStatistics(bob, bobPoll),
);
