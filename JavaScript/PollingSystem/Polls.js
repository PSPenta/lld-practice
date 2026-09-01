const { Poll } = require('./Poll');

class Polls {
  polls = [];

  static addPoll(poll) {
    if (!(poll instanceof Poll)) {
      throw new Error('Invalid poll');
    }

    if (!this.polls) {
      this.polls = [];
    }

    this.polls.push(poll);
  }

  static getActivePolls() {
    return this.polls.filter(poll => poll.validTill > Date.now());
  }

  static getCompletedPolls(adminId) {
    return this.polls.filter(poll => poll.createdBy == adminId && poll.validTill < Date.now());
  }
}

module.exports = {Polls};