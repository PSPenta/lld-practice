package main

import "fmt"

func main() {
	user1 := NewUser("XfH4I@example.com")
	user2 := NewUser("4fNl4@example.com")
	user3 := NewUser("7Eo5o@example.com")

	balanceSheet := NewBalanceSheet()

	split1, _ := NewSplit(user1.Email, 0, 0)
	split2, _ := NewSplit(user2.Email, 0, 0)
	split3, _ := NewSplit(user3.Email, 0, 0)
	expense, _ := CreateExpense("Equal", user1.Email, 300, []Split{*split1, *split2, *split3})
	expense.Validate()
	expense.Apply(balanceSheet)
	fmt.Println(balanceSheet.GetBalances())

	split4, _ := NewSplit(user1.Email, 0, 10)
	split5, _ := NewSplit(user2.Email, 0, 20)
	split6, _ := NewSplit(user3.Email, 0, 70)
	expense2, _ := CreateExpense("Percentage", user2.Email, 150, []Split{*split4, *split5, *split6})
	expense2.Validate()
	expense2.Apply(balanceSheet)
	fmt.Println(balanceSheet.GetBalances())

	split7, _ := NewSplit(user1.Email, 225, 0)
	split8, _ := NewSplit(user2.Email, 390, 0)
	split9, _ := NewSplit(user3.Email, 680, 0)
	expense3, _ := CreateExpense("Exact", user3.Email, 1295, []Split{*split7, *split8, *split9})
	expense3.Validate()
	expense3.Apply(balanceSheet)
	fmt.Println(balanceSheet.GetBalances())

	fmt.Println(expense)
	fmt.Println(expense2)
	fmt.Println(expense3)
	fmt.Println(balanceSheet.GetBalances())
}
