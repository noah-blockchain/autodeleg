package helpers

import (
	"math/big"
	"testing"
)

func TestStringToBigInt(t *testing.T) {
	value := "100"
	expected := big.NewInt(100)
	got := StringToBigInt(value)
	if got.Cmp(expected) != 0 {
		t.Errorf("got %v wanted %v", got, expected)
	}

	value = ""
	expected = big.NewInt(0)
	got = StringToBigInt(value)
	if got.Cmp(expected) != 0 {
		t.Errorf("got %v wanted %v", got, expected)
	}
}
