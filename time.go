package awips

import (
	"strings"
	"time"
)

// Common timezones used in AWIPS products.
var Timezones = map[string]*time.Location{
	// GMT/UTC/Zulu
	"GMT": time.FixedZone("GMT", 0*60*60),
	"UTC": time.FixedZone("UTC", 0*60*60),
	// Atlantic
	"ADT": time.FixedZone("ADT", -3*60*60),
	"AST": time.FixedZone("AST", -4*60*60),
	// Easter
	"EST": time.FixedZone("EST", -5*60*60),
	"EDT": time.FixedZone("EDT", -4*60*60),
	// Central
	"CST": time.FixedZone("CST", -6*60*60),
	"CDT": time.FixedZone("CDT", -5*60*60),
	// Mountain
	"MST": time.FixedZone("MST", -7*60*60),
	"MDT": time.FixedZone("MDT", -6*60*60),
	// Pacific
	"PST": time.FixedZone("PST", -8*60*60),
	"PDT": time.FixedZone("PDT", -7*60*60),
	// Alaska
	"AKST": time.FixedZone("AKST", -9*60*60),
	"AKDT": time.FixedZone("AKDT", -8*60*60),
	// Hawaii
	"HST": time.FixedZone("HST", -10*60*60),
	// Samoa
	"SST": time.FixedZone("SST", -11*60*60),
	// Chamorro/Guam
	"CHST": time.FixedZone("CHST", 10*60*60),
}

func GetTimezone(tzs string) *time.Location {
	tzs = strings.ToUpper(tzs)

	for key, timezone := range Timezones {
		if key == tzs {
			return timezone
		}
	}

	return nil
}
