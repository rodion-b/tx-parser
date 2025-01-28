package utils

import (
	"regexp"
	"strings"
)

func IsValidEthereumAddress(address string) bool {
	if len(address) != 42 || !strings.HasPrefix(address, "0x") {
		return false
	}

	// Regular expression to check if the rest of the string is valid hex
	regex := regexp.MustCompile("^(0x)[0-9a-fA-F]{40}$")
	return regex.MatchString(address)
}
