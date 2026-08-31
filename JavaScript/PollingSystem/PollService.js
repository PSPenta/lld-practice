const { User } = require('./User');
const { Poll } = require('./Poll');
const { Vote } = require('./Vote');
const { UserRepository } = require('./UserRepository');
const { PollRepository } = require('./PollRepository');
const { VoteRepository } = require('./VoteRepository');

/**
 * Application / use-case layer for polling.
 * Entities stay thin; authorization + orchestration live here.
 */
class PollService {
  createUser(email) {
    return UserRepository.create(email);
  }

  createPoll(creator, question, options, durationMs) {
    if (!(creator instanceof User)) {
      throw new Error('Invalid creator!');
    }

    const poll = new Poll(
      PollRepository.createId(),
      question,
      options,
      creator.id,
      durationMs,
    );
    PollRepository.add(poll);
    return poll;
  }

  /** Only the poll creator can invite voters. */
  assignVoter(creator, poll, voter) {
    if (!(creator instanceof User) || !(poll instanceof Poll) || !(voter instanceof User)) {
      throw new Error('Invalid creator, poll, or voter!');
    }

    if (!poll.isCreator(creator.id)) {
      throw new Error('Only the poll creator can assign voters!');
    }

    if (poll.isExpired()) {
      throw new Error('Poll has been expired!');
    }

    poll.assignVoter(voter.id);
  }

  /**
   * Cast a vote. Users may vote on others' polls, never on their own.
   */
  submitVote(voter, poll, option) {
    if (!(voter instanceof User) || !(poll instanceof Poll) || !option) {
      throw new Error('Invalid voter, poll, or option!');
    }

    if (poll.isCreator(voter.id)) {
      throw new Error('You cannot vote on your own poll!');
    }

    if (poll.isExpired()) {
      throw new Error('Poll has been expired!');
    }

    if (!poll.options.includes(option)) {
      throw new Error('Invalid option!');
    }

    if (!poll.isAssigned(voter.id)) {
      throw new Error('You are not assigned to this poll!');
    }

    VoteRepository.add(new Vote(poll.id, option, voter.id));
  }

  /** Only the poll creator can view statistics. */
  getStatistics(creator, poll) {
    if (!(creator instanceof User) || !(poll instanceof Poll)) {
      throw new Error('Invalid creator or poll!');
    }

    if (!poll.isCreator(creator.id)) {
      throw new Error('Only the poll creator can view statistics!');
    }

    return VoteRepository.getStatistics(poll);
  }

  getActivePolls() {
    return PollRepository.getActive();
  }

  getCompletedPolls(creatorId) {
    return PollRepository.getCompletedByCreator(creatorId);
  }
}

module.exports = { PollService };
