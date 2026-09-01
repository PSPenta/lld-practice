package pollingservice

import (
	"fmt"
	"time"

	"lld-practice/pollingsystem2-go/models"
	"lld-practice/pollingsystem2-go/repositories"
)

// PollingService orchestrates use-cases (mirrors JavaScript/PollingSystem2/PollingService/).
type PollingService struct {
	users *repositories.UserRepository
	polls *repositories.PollRepository
	votes *repositories.VoteRepository
}

func New(users *repositories.UserRepository, polls *repositories.PollRepository, votes *repositories.VoteRepository) *PollingService {
	return &PollingService{users: users, polls: polls, votes: votes}
}

func (s *PollingService) CreateUser(email string) (*models.User, error) {
	return s.users.Create(email)
}

func (s *PollingService) CreatePoll(creator *models.User, question string, options []string, isPrivate, isClosed bool, duration time.Duration) (*models.Poll, error) {
	if creator == nil {
		return nil, fmt.Errorf("invalid creator!")
	}
	poll, err := models.NewPoll(s.polls.NextID(), question, options, creator.ID, isPrivate, isClosed, duration)
	if err != nil {
		return nil, err
	}
	if err := s.polls.Add(poll); err != nil {
		return nil, err
	}
	return poll, nil
}

func (s *PollingService) GetPoll(pollID int) *models.Poll {
	return s.polls.GetByID(pollID)
}

func (s *PollingService) UpdatePoll(poll *models.Poll) error {
	if poll == nil {
		return fmt.Errorf("invalid poll!")
	}
	return s.polls.Update(poll)
}

func (s *PollingService) AssignVoter(creator *models.User, poll *models.Poll, voter *models.User) error {
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

func (s *PollingService) SubmitVote(voter *models.User, poll *models.Poll, option string) error {
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
	vote, err := models.NewVote(poll.ID, option, voter.ID)
	if err != nil {
		return err
	}
	return s.votes.Add(vote)
}

func (s *PollingService) GetStatistics(creator *models.User, poll *models.Poll) (map[string]any, error) {
	if creator == nil || poll == nil {
		return nil, fmt.Errorf("invalid creator or poll!")
	}
	if !poll.IsCreator(creator.ID) {
		return nil, fmt.Errorf("only the poll creator can view statistics!")
	}
	return s.votes.GetStatistics(poll), nil
}

func (s *PollingService) GetActivePolls() []*models.Poll {
	return s.polls.GetActive(time.Now())
}

func (s *PollingService) GetCompletedPolls(creatorID int) []*models.Poll {
	return s.polls.GetCompletedByCreator(creatorID, time.Now())
}
