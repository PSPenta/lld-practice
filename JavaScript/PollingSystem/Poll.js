class Poll {
  id = 0;
  question = '';
  options = [];
  scheduledTime = 0;
  validTill = 0;
  createdBy = 0;
  /** @type {number[]} */
  assignedUsers = [];

  /**
   * @param {number} id
   * @param {string} question
   * @param {string[]} options
   * @param {number} createdBy
   * @param {number} [durationMs=24 * 60 * 60 * 1000]
   */
  constructor(id, question, options, createdBy, durationMs = 24 * 60 * 60 * 1000) {
    if (!id || !question || !options || options.length <= 1 || !createdBy) {
      throw new Error('Invalid poll parameters!');
    }

    const uniqueOptions = [...new Set(options)];
    if (uniqueOptions.length !== options.length) {
      throw new Error('Poll options must be unique!');
    }

    const now = Date.now();
    this.id = id;
    this.question = question;
    this.options = options;
    this.createdBy = createdBy;
    this.scheduledTime = now;
    this.validTill = now + durationMs;
  }

  /** @param {number} userId */
  assignVoter(userId) {
    if (!userId) {
      throw new Error('Invalid user!');
    }

    if (userId === this.createdBy) {
      throw new Error('You cannot assign yourself to your own poll!');
    }

    if (this.assignedUsers.includes(userId)) {
      throw new Error('User already assigned to this poll!');
    }

    this.assignedUsers.push(userId);
  }

  isCreator(userId) {
    return this.createdBy === userId;
  }

  isAssigned(userId) {
    return this.assignedUsers.includes(userId);
  }

  isExpired(now = Date.now()) {
    return this.validTill < now;
  }
}

module.exports = { Poll };
