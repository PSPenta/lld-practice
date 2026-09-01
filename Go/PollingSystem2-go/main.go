package main

import (
	"fmt"

	"lld-practice/pollingsystem2-go/pollingservice"
	"lld-practice/pollingsystem2-go/repositories"
)

func main() {
	svc := pollingservice.New(
		repositories.NewUserRepository(),
		repositories.NewPollRepository(),
		repositories.NewVoteRepository(),
	)

	alice, _ := svc.CreateUser("alice@example.com")
	bob, _ := svc.CreateUser("bob@example.com")
	jane, _ := svc.CreateUser("jane@example.com")
	jim, _ := svc.CreateUser("jim@example.com")
	kate, _ := svc.CreateUser("kate@example.com")

	alicePoll, err := svc.CreatePoll(alice, "Whats capital of India?", []string{"Delhi", "Mumbai"}, true, false, 0)
	if err != nil {
		panic(err)
	}

	if err := svc.AssignVoter(alice, alicePoll, alice); err != nil {
		fmt.Println("Expected error TestCase 1:", err.Error())
	}
	if err := svc.SubmitVote(alice, alicePoll, "Delhi"); err != nil {
		fmt.Println("Expected error TestCase 2:", err.Error())
	}
	if err := svc.SubmitVote(bob, alicePoll, "Mumbai"); err != nil {
		fmt.Println("Expected error TestCase 3:", err.Error())
	}

	_ = svc.AssignVoter(alice, alicePoll, bob)
	_ = svc.SubmitVote(bob, alicePoll, "Mumbai")
	_ = svc.AssignVoter(alice, alicePoll, jane)
	_ = svc.SubmitVote(jane, alicePoll, "Delhi")
	_ = svc.AssignVoter(alice, alicePoll, jim)
	_ = svc.SubmitVote(jim, alicePoll, "Delhi")

	alicePoll.IsClosed = true
	_ = svc.UpdatePoll(alicePoll)

	if err := svc.AssignVoter(alice, alicePoll, kate); err != nil {
		fmt.Println("Expected error TestCase 6:", err.Error())
	}
	if err := svc.SubmitVote(jim, alicePoll, "Delhi"); err != nil {
		fmt.Println("Expected error TestCase 7:", err.Error())
	}

	stats, _ := svc.GetStatistics(alice, alicePoll)
	fmt.Printf("Poll statistics for %d: %+v\n", alicePoll.ID, stats)

	bobPoll, _ := svc.CreatePoll(bob, "Best language?", []string{"JavaScript", "Go"}, false, false, 0)
	_ = svc.SubmitVote(alice, bobPoll, "JavaScript")
	bobStats, _ := svc.GetStatistics(bob, bobPoll)
	fmt.Printf("Poll statistics for %d: %+v\n", bobPoll.ID, bobStats)
}
