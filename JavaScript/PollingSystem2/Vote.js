class Vote {
  pollId = 0;
  option = '';
  userId = 0;
  submittedAt = 0;

  constructor(pollId, option, userId) {
    if (!pollId || !option || !userId) {
      throw new Error('Invalid vote parameters!');
    }

    this.pollId = pollId;
    this.option = option;
    this.userId = userId;
    this.submittedAt = Date.now();
  }
}

module.exports = { Vote };
