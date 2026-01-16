package awips

import (
	"errors"
	"regexp"
	"strings"
)

// AWIPS Header
type AWIPS struct {
	Original string `json:"original"`
	Product  string `json:"product"` // Product category
	NWSLI    string `json:"wfo"`     // NWS Location Identifier
}

var ErrCouldNotFindAWIPS = errors.New("could not find AWIPS header")

// This regular expression is derived from NWS directive 10-1701 section 4.1.3.
// The directive states that the AWIPS header is a 4 to 6 character string on its own line.
// The string can contain letters and numbers.
const AWIPSRegexp = `(?m:^[A-Z0-9]{3}[A-Z0-9 ]{1,3}[\n\r])`

// Returns the AWIPS header from the given text.
// If no header is found, the string is empty.
func FindAWIPS(text string) string {
	awipsRegex := regexp.MustCompile(AWIPSRegexp)
	return strings.TrimSpace(awipsRegex.FindString(text))
}

// Checks if the given text contains an AWIPS header.
func HasAWIPS(text string) bool {
	return FindAWIPS(text) != ""
}

// Parses the AWIPS header from the given text.
func ParseAWIPS(text string) (AWIPS, error) {
	// Find the AWIPS header
	original := FindAWIPS(text)
	if original == "" {
		return AWIPS{}, ErrCouldNotFindAWIPS
	}

	// Product is the first three characters
	product := strings.TrimSpace(original[0:3])
	// Issuing office is the final three products
	wfo := strings.TrimSpace(original[3:])

	return AWIPS{
		Original: original,
		Product:  product,
		NWSLI:    wfo,
	}, nil
}
