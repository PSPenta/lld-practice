const {Poll} = require('./Poll');
const {Results} = require('./Results');

class Admin {
  id = 0;
  constructor(id) {
    this.id = id;
  }

  createPoll(question, options) {
    const poll = new Poll(question, options, this.id);

    return poll;
  }

  showStatistics(poll) {
    const statistics = Results.getStatistics(poll);
    console.log(`Poll statistics for ${poll.id}:`, statistics);
  }
}

module.exports = {Admin};