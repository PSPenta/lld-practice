const { PollingService } = require('./PollingService');

const pollService = new PollingService();

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

// Test case 1: Assign voter to self
try {
  pollService.assignVoter(alice, alicePoll, alice);
} catch (err) {
  console.log('Expected error TestCase 1:', err.message);
}

// Test case 2: Submit vote to self-created poll
try {
  pollService.submitVote(alice, alicePoll, 'Delhi');
} catch (err) {
  console.log('Expected error TestCase 2:', err.message);
}

// Test case 3: Submit vote to another user's private poll
try {
  pollService.submitVote(bob, alicePoll, 'Mumbai');
} catch (err) {
  console.log('Expected error TestCase 3:', err.message);
}

// Test case 4: Assign voter to another user
pollService.assignVoter(alice, alicePoll, bob);
pollService.submitVote(bob, alicePoll, 'Mumbai');

pollService.assignVoter(alice, alicePoll, jane);
pollService.submitVote(jane, alicePoll, 'Delhi');

pollService.assignVoter(alice, alicePoll, jim);
pollService.submitVote(jim, alicePoll, 'Delhi');

// Test case 5: Close poll
alicePoll.isClosed = true;
pollService.updatePoll(alicePoll);

// Test case 6: Assign voter to closed poll
try {
  pollService.assignVoter(alice, alicePoll, kate);
} catch (err) {
  console.log('Expected error TestCase 6:', err.message);
}

// Test case 7: Submit vote to closed poll
try {
  pollService.submitVote(jim, alicePoll, 'Delhi');
} catch (err) {
  console.log('Expected error TestCase 7:', err.message);
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
