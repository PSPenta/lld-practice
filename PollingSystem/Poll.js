class Poll {
  question = '';
  options = {};
  scheduledTime = Date();
  validTill = Date();
  createBy = 0;

  constructor(question, options, createdBy) {
    if (!question || question == '' || !options || options.length <= 1 || !createdBy) {
      throw new Error('Invalid poll parameters!');
    }
    this.question = question;
    this.options = options;
    this.createBy = createdBy;
  }
}

module.exports = {Poll};