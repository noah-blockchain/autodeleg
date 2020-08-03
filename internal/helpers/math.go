package helpers

import (
	"math/big"
)

func StringToBigInt(string string) *big.Int {
	bInt, err := new(big.Int).SetString(string, 10)
	CheckErrBool(err)

	return bInt
}
