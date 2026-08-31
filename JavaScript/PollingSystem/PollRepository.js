const { Poll } = require('./Poll');

/** In-memory stand-in for a poll repository. */
class PollRepository {
  static polls = [];
  static nextId = 1;

  static createId() {
    return this.nextId++;
  }

  static add(poll) {
    if (!(poll instanceof Poll)) {
      throw new Error('Invalid poll');
    }

    if (this.polls.some((p) => p.id === poll.id)) {
      throw new Error('Poll already exists!');
    }

    this.polls.push(poll);
  }

  static getById(id) {
    return this.polls.find((p) => p.id === id);
  }

  static getActive(now = Date.now()) {
    return this.polls.filter((poll) => !poll.isExpired(now));
  }

  static getCompletedByCreator(creatorId, now = Date.now()) {
    return this.polls.filter(
      (poll) => poll.createdBy === creatorId && poll.isExpired(now),
    );
  }
}

module.exports = { PollRepository };
