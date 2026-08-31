class Result {
  pollId = 0;
  createdBy = 0;
  option = '';
  userId = 0;
  submittedAt = new Date();

  constructor(pollId, option, userId) {
    this.pollId = pollId;
    this.option = option;
    this.userId = userId;
    this.submittedAt = Date.now();
  }
}

module.exports = {Result};