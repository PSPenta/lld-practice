const { Vote } = require('../models/Vote');

/** In-memory stand-in for a vote repository. */
class VoteRepository {
  static votes = [];

  static add(vote) {
    if (!(vote instanceof Vote)) {
      throw new Error('Invalid vote');
    }

    const alreadyVoted = this.votes.some(
      (v) => v.pollId === vote.pollId && v.userId === vote.userId,
    );
    if (alreadyVoted) {
      throw new Error('User has already voted on this poll!');
    }

    this.votes.push(vote);
  }

  static getByPollId(pollId) {
    return this.votes.filter((v) => v.pollId === pollId);
  }

  static getStatistics(poll) {
    const votes = this.getByPollId(poll.id);
    const total = votes.length;

    const ratings = {};
    for (const vote of votes) {
      ratings[vote.option] = (ratings[vote.option] || 0) + 1;
    }

    const statistics = {};
    for (const option of poll.options) {
      const count = ratings[option] || 0;
      statistics[option] = total === 0 ? 0 : count / total;
    }

    return {
      question: poll.question,
      totalVotes: total,
      counts: Object.fromEntries(
        poll.options.map((option) => [option, ratings[option] || 0]),
      ),
      statistics,
    };
  }
}

module.exports = { VoteRepository };
