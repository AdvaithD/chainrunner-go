package util

import (
	"math/big"
)

// Removes leading zeroes from a string
func RemoveLeadingZeros(inputString string) string {
	for index := 0; index < len(inputString); index++ {
		// Just return the string as soon as first non zero value is detected
		if string(inputString[index]) != "0" {
			return inputString[index:]
		}
	}
	return "" // Value is only zeros
}

// Decode Reserve from storage slot
func DeriveReservesFromSlot(slot string) (*big.Int, *big.Int) {
	// Yes these names are right, for some reason they are stored in reverse order
	// The below two lines use uint256.Int (TODO: future optimization)

	// reserve1, _ := uint256.FromHex("0x" + RemoveLeadingZeros(slot[10:38]))
	// reserve0, _ := uint256.FromHex("0x" + RemoveLeadingZeros(slot[38:66]))

	// The below uses big.Int
	reserve1, _ := new(big.Int).SetString("0x"+RemoveLeadingZeros(slot[10:38]), 0)

	reserve0, _ := new(big.Int).SetString("0x"+RemoveLeadingZeros(slot[38:66]), 0)

	return reserve0, reserve1
}
