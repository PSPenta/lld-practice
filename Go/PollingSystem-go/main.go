package main

func main() {
	admin := NewAdmin(1)
	poll, _ := admin.CreatePoll("Whats capital of India?", []string{"Delhi", "Mumbai"})
	AddPoll(poll)

	user := NewUser(2)
	user.SubmitPoll(poll, "Mumbai")

	user2 := NewUser(3)
	user2.SubmitPoll(poll, "Delhi")

	user3 := NewUser(3)
	user3.SubmitPoll(poll, "Delhi")

	user4 := NewUser(3)
	user4.SubmitPoll(poll, "Delhi")

	admin.ShowStatistics(poll)
}
