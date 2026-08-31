package main

import (
	"fmt"
	"time"
)

func main() {
	svc := NewPollService(NewUserRepository(), NewPollRepository(), NewVoteRepository())

	alice, _ := svc.CreateUser("alice@example.com")
	bob, _ := svc.CreateUser("bob@example.com")
	jane, _ := svc.CreateUser("jane@example.com")
	jim, _ := svc.CreateUser("jim@example.com")

	alicePoll, err := svc.CreatePoll(alice, "Whats capital of India?", []string{"Delhi", "Mumbai"}, 24*time.Hour)
	if err != nil {
		panic(err)
	}

	if err := svc.AssignVoter(alice, alicePoll, alice); err != nil {
		fmt.Println("Expected error:", err.Error())
	}
	if err := svc.SubmitVote(alice, alicePoll, "Delhi"); err != nil {
		fmt.Println("Expected error:", err.Error())
	}
	if err := svc.SubmitVote(bob, alicePoll, "Mumbai"); err != nil {
		fmt.Println("Expected error:", err.Error())
	}

	_ = svc.AssignVoter(alice, alicePoll, bob)
	_ = svc.SubmitVote(bob, alicePoll, "Mumbai")
	_ = svc.AssignVoter(alice, alicePoll, jane)
	_ = svc.SubmitVote(jane, alicePoll, "Delhi")
	_ = svc.AssignVoter(alice, alicePoll, jim)
	_ = svc.SubmitVote(jim, alicePoll, "Delhi")

	if err := svc.SubmitVote(jim, alicePoll, "Delhi"); err != nil {
		fmt.Println("Expected error:", err.Error())
	}

	stats, _ := svc.GetStatistics(alice, alicePoll)
	fmt.Printf("Poll statistics for %d: %+v\n", alicePoll.ID, stats)

	bobPoll, _ := svc.CreatePoll(bob, "Best language?", []string{"JavaScript", "Go"}, 24*time.Hour)
	_ = svc.AssignVoter(bob, bobPoll, alice)
	_ = svc.SubmitVote(alice, bobPoll, "JavaScript")
	bobStats, _ := svc.GetStatistics(bob, bobPoll)
	fmt.Printf("Poll statistics for %d: %+v\n", bobPoll.ID, bobStats)
}
