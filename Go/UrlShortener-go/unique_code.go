package main

const (
	base62        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultLength = 6
	secret        = 9765439
)

func generateShortCode(counterStore map[string]int) string {
	counter := counterStore["UNIQUE_COUNTER"] + 1
	counterStore["UNIQUE_COUNTER"] = counter

	temp := counter ^ secret
	shortCode := make([]byte, defaultLength)
	for i := defaultLength - 1; i >= 0; i-- {
		shortCode[i] = base62[temp%62]
		temp /= 62
	}
	return string(shortCode)
}
