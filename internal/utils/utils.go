/*
Copyright © 2024 Patrick X. Gray pxgray@proton.me
*/
package utils

import (
	"math/big"

	"github.com/shopspring/decimal"
)

// Convert a string-based decimal amount to wei
func ToWei(v string, decimals int) *big.Int {
	amount, _ := decimal.NewFromString(v)

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}
