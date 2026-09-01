const { Result } = require('./Result');

class Results {
  results = [];

  static submission(result) {
    if (!(result instanceof Result)) {
      throw new Error('Invalid submission');
    }

    if (!this.results) {
      this.results = [];
    }

    this.results.push(result);
  }

  static getStatistics(poll) {
    let results = this.results.filter(res => res.pollId == poll.id);

    let ratings = {};
    for (let res of results) {
      if (!ratings[res.option]) {
        ratings[res.option] = 0;
      }

      ratings[res.option]++;
    }

    let res = {};
    for (let option of poll.options) {
      res[option] = ratings[option] / results.length;
    }

    return {
      question: poll.question,
      statistics: res
    };
  }
}

module.exports = {Results}