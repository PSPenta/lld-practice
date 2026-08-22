package main

import "fmt"

func main() {
	engine := NewSearchEngine()
	engine.AddDocument("1", "Apple launches new MacBook Pro with M4 chip")
	engine.AddDocument("2", "Samsung introduces AI-powered Galaxy devices")
	engine.AddDocument("3", "Apple releases major update for iOS 19")

	fmt.Println("🔍 Search: 'apple m4'")
	fmt.Println(engine.Search("apple m4", 10))

	fmt.Println("🔍 Search: 'apple 19'")
	fmt.Println(engine.Search("apple 19", 10))

	fmt.Println("💡 Autocomplete: 'mac'")
	fmt.Println(engine.Autocomplete("mac", 5))
}
