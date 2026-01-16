package awips

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type WMO struct {
	Original string    `json:"original"`
	Datatype string    `json:"datatype"`
	Office   string    `json:"office"`
	Issued   time.Time `json:"issued"` // Only day, hour, minute
	BBB      string    `json:"bbb"`
}

const WMORegexp = `([A-Z]{4}[0-9]{2})\s([A-Z]{4})\s([0-9]{6})( [A-Z]{3})?`

// Attempts to find the WMO line in the provided text.
// Empty string is returned if no match.
func FindWMO(text string) string {
	wmoRegexp := regexp.MustCompile(WMORegexp)
	return wmoRegexp.FindString(text)
}

// Parse the WMO line into information we can use.
func ParseWMO(text string) (WMO, error) {
	// Find the WMO line
	wmoRegexp := regexp.MustCompile(WMORegexp)
	original := wmoRegexp.FindString(text)
	if original == "" {
		return WMO{}, errors.New("could not find WMO line")
	}

	// Segment the line
	segments := strings.Split(original, " ")

	// Time layout (ddhhmm)
	layout := "021504"

	// Issued day & time
	t, err := time.Parse(layout, segments[2])
	if err != nil {
		return WMO{}, errors.New("could not parse WMO issued datetime: " + err.Error())
	}

	// bbb if any exists
	bbb := ""
	if len(segments) > 3 {
		bbb = segments[3]
	}

	return WMO{
		Original: original,
		Datatype: segments[0],
		Office:   segments[1],
		Issued:   t,
		BBB:      bbb,
	}, nil
}

// Checks if the provided text contains a WMO line.
func HasWMO(text string) bool {
	original := FindWMO(text)
	return original != ""
}
