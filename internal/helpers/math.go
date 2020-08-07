package helpers

import (
	"math/big"
)

func StringToBigInt(string string) *big.Int {
	if string == "" {
		return big.NewInt(0)
	}
	bInt, err := new(big.Int).SetString(string, 10)
	CheckErrBool(err)
	return bInt
}
