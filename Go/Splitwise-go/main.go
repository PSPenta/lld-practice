package main

import "fmt"

func main() {
	service := NewSplitwiseService()

	user1, _ := service.AddUser("John Doe", "john.doe@example.com")
	if _, err := service.AddUser("John Doe", "john.doe@example.com"); err != nil {
		fmt.Println("Error adding user:", err)
	}
	user2, _ := service.AddUser("Jane Doe", "jane.doe@example.com")
	user3, _ := service.AddUser("Jim Doe", "jim.doe@example.com")

	_, _ = service.AddExpense(AddExpenseInput{
		Type:   "Equal",
		PaidBy: user1.ID,
		Amount: 300,
		Splits: []*Split{
			NewEqualSplit(user1.ID),
			NewEqualSplit(user2.ID),
			NewEqualSplit(user3.ID),
		},
	})
	fmt.Println(service.GetPairwiseBalances())

	_, _ = service.AddExpense(AddExpenseInput{
		Type:   "percentage",
		PaidBy: user2.ID,
		Amount: 150.50,
		Splits: []*Split{
			mustPct(user1.ID, 10),
			mustPct(user2.ID, 20),
			mustPct(user3.ID, 70),
		},
	})
	fmt.Println(service.GetPairwiseBalances())

	_ = service.SettleUp(user2.ID, user1.ID, nil)
	fmt.Println(service.GetPairwiseBalances())

	_ = service.SettleUp(user3.ID, user2.ID, nil)
	fmt.Println(service.GetPairwiseBalances())

	partial := 99.0
	_ = service.SettleUp(user3.ID, user1.ID, &partial)
	fmt.Println(service.GetPairwiseBalances())
}

func mustPct(userID, pct int) *Split {
	s, err := NewPercentageSplit(userID, pct)
	if err != nil {
		panic(err)
	}
	return s
}
