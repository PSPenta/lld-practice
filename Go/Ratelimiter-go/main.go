package main

import "fmt"

func main() {
	limiter := NewFixedWindowRatelimiter(1000, 5)
	fmt.Println("Fixed window allowed:", limiter.IsRequestAllowed("user1", "/api"))

	tbLimiter := NewTokenBucketLimiter()
	fmt.Println("Token bucket allowed:", tbLimiter.AllowRequest("127.0.0.1"))

	lbLimiter := NewLeakyBucketLimiter()
	fmt.Println("Leaky bucket allowed:", lbLimiter.AllowRequest("127.0.0.1"))

	demoConcurrencyThrottler()
}
