package main

import (
	"fmt"
	"time"
)

type PollService struct {
	users *UserRepository
	polls *PollRepository
	votes *VoteRepository
}

func NewPollService(users *UserRepository, polls *PollRepository, votes *VoteRepository) *PollService {
	return &PollService{users: users, polls: polls, votes: votes}
}

func (s *PollService) CreateUser(email string) (*User, error) {
	return s.users.Create(email)
}

func (s *PollService) CreatePoll(creator *User, question string, options []string, isPrivate, isClosed bool, duration time.Duration) (*Poll, error) {
	if creator == nil {
		return nil, fmt.Errorf("invalid creator!")
	}
	poll, err := NewPoll(s.polls.NextID(), question, options, creator.ID, isPrivate, isClosed, duration)
	if err != nil {
		return nil, err
	}
	if err := s.polls.Add(poll); err != nil {
		return nil, err
	}
	return poll, nil
}

func (s *PollService) GetPoll(pollID int) *Poll {
	return s.polls.GetByID(pollID)
}

func (s *PollService) UpdatePoll(poll *Poll) error {
	if poll == nil {
		return fmt.Errorf("invalid poll!")
	}
	return s.polls.Update(poll)
}

func (s *PollService) AssignVoter(creator *User, poll *Poll, voter *User) error {
	if creator == nil || poll == nil || voter == nil {
		return fmt.Errorf("invalid creator, poll, or voter!")
	}
	if !poll.IsCreator(creator.ID) {
		return fmt.Errorf("only the poll creator can assign voters!")
	}
	if poll.IsCompleted(time.Now()) {
		return fmt.Errorf("poll has been expired!")
	}
	if !poll.IsPrivate {
		return fmt.Errorf("cannot assign voters to public polls!")
	}
	if poll.IsClosed {
		return fmt.Errorf("poll has been completed!")
	}
	return poll.AssignVoter(voter.ID)
}

func (s *PollService) SubmitVote(voter *User, poll *Poll, option string) error {
	if voter == nil || poll == nil || option == "" {
		return fmt.Errorf("invalid voter, poll, or option!")
	}
	if poll.IsCreator(voter.ID) {
		return fmt.Errorf("you cannot vote on your own poll!")
	}
	if poll.IsCompleted(time.Now()) {
		return fmt.Errorf("poll has been expired!")
	}
	valid := false
	for _, o := range poll.Options {
		if o == option {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid option!")
	}
	if poll.IsPrivate && !poll.IsAssigned(voter.ID) {
		return fmt.Errorf("you are not assigned to this poll!")
	}
	if poll.IsClosed {
		return fmt.Errorf("poll has been completed!")
	}
	vote, err := NewVote(poll.ID, option, voter.ID)
	if err != nil {
		return err
	}
	return s.votes.Add(vote)
}

func (s *PollService) GetStatistics(creator *User, poll *Poll) (map[string]any, error) {
	if creator == nil || poll == nil {
		return nil, fmt.Errorf("invalid creator or poll!")
	}
	if !poll.IsCreator(creator.ID) {
		return nil, fmt.Errorf("only the poll creator can view statistics!")
	}
	return s.votes.GetStatistics(poll), nil
}

func (s *PollService) GetActivePolls() []*Poll {
	return s.polls.GetActive(time.Now())
}

func (s *PollService) GetCompletedPolls(creatorID int) []*Poll {
	return s.polls.GetCompletedByCreator(creatorID, time.Now())
}
