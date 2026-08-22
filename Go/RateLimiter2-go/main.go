package main

import "fmt"

func main() {
	rateLimiter := NewRateLimiter(NewTokenBucketStrategy(5, 1))
	rateLimiter.IsAllowed("127.0.0.1")

	rateLimiter2 := NewRateLimiter(NewFixedWindowCounter(5, 1000))
	rateLimiter2.IsAllowed("127.0.0.1")

	rateLimiter3 := NewRateLimiter(NewSlidingWindowLog(5, 1000))
	rateLimiter3.IsAllowed("127.0.0.1")

	rateLimiter4 := NewRateLimiter(NewLeakyBucketStrategy(5, 1))
	rateLimiter4.IsAllowed("127.0.0.1")

	rateLimiter5 := NewRateLimiter(NewSlidingWindowCounter(5, 1000))
	rateLimiter5.IsAllowed("127.0.0.1")

	fmt.Println(rateLimiter.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter2.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter3.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter4.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter5.IsAllowed("127.0.0.1"))

	fmt.Println(rateLimiter.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter2.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter3.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter4.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter5.IsAllowed("127.0.0.1"))

	fmt.Println(rateLimiter.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter2.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter3.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter4.IsAllowed("127.0.0.1"))
	fmt.Println(rateLimiter5.IsAllowed("127.0.0.1"))
}
