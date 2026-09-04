package main

import (
	"fmt"
	"math"
)

// ToAmount converts rupees to integer paise. Ledger math stays in paise.
func ToAmount(rupees float64) (int64, error) {
	if math.IsNaN(rupees) || math.IsInf(rupees, 0) {
		return 0, fmt.Errorf("invalid amount")
	}
	return int64(math.Round(rupees * 100)), nil
}

func FromAmount(paise int64) float64 {
	return float64(paise) / 100
}
