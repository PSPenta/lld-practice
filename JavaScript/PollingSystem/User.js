const { Poll } = require('./Poll');
const { Result } = require('./Result');
const { Results } = require('./Results');

class User {
  id = 0;

  constructor(id) {
    this.id = id;
  }

  submitPoll(poll, option) {
    if (!(poll instanceof Poll) || !option) {
      throw new Error('Invalid poll or option');
    }

    if (poll.validTill < Date.now()) {
      throw new Error('Poll has been expired!');
    }

    if (!poll.options.includes(option)) {
      throw new Error('Invalid option!');
    }

    Results.submission(new Result(poll.id, option, this.id));
  }
}

module.exports = {User};