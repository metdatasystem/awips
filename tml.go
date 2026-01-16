package awips

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/twpayne/go-geom"
)

// The rules of TIME...MOT...LOC information in text products is defined in NWS directive 10-1701 section 5.7.

const TMLRegexp = `(?m:^(TIME\.\.\.MOT\.\.\.LOC)([A-Za-z0-9 ]+\n{0,1})*)`

type TML struct {
	Original    string           `json:"original"`
	Time        time.Time        `json:"time"`
	Direction   int              `json:"direction"`
	Speed       int              `json:"speed"`
	SpeedString string           `json:"speedString"`
	Locations   *geom.MultiPoint `json:"location"`
}

// Find the TIME...MOT...LOC information in the text.
func FindTML(text string) string {
	tmlRegexp := regexp.MustCompile(TMLRegexp)
	return strings.TrimSpace(tmlRegexp.FindString(text))
}

// Attempts to find the TIME...MOT...LOC information in the text and parse it.
// If the string is not found, it returns nil.
// If the string is found but cannot be parsed, it returns an error.
func ParseTML(text string, issued time.Time) (*TML, error) {

	original := FindTML(text)
	if original == "" {
		return nil, nil
	}

	trimRegexp := regexp.MustCompile(`[\s\n]+`)
	original = trimRegexp.ReplaceAllString(original, " ")

	// Split the string into segments. Segments are separated by spaces.
	segments := strings.Split(original, " ")[1:]
	if len(segments) == 0 {
		return nil, errors.New("tml segments is 0")
	}

	// Parse the time
	parsedTime, err := time.Parse(("1504Z"), segments[0])
	if err != nil {
		return nil, errors.New("could not parse TML time: " + err.Error())
	}

	time := time.Date(issued.Year(), issued.Month(), issued.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Now().UTC().Location())

	// Parse the direction
	direction, err := strconv.Atoi(segments[1][:3])
	if err != nil {
		return nil, errors.New("could not parse direction in TML: " + err.Error())
	}

	numberRegexp := regexp.MustCompile("[0-9]+")

	// Parse the speed
	// The directive says the speed is given as knots. However, some products use miles per hour.
	speedString := segments[2]
	speed, err := strconv.Atoi(numberRegexp.FindString(speedString))
	if err != nil {
		return nil, errors.New("could not parse speed in TML: " + err.Error())
	}

	points := geom.NewMultiPoint(geom.XY)

	// Parse the location(s).
	// TML information can contain multiple pairs of latitude and longitude information whether the product is talking about an event that is isolated to a single point or the event is occurring along a line
	for i := 3; i < len(segments)-1; i += 2 {
		lat, err := Parse45Digit(segments[i])
		if err != nil {
			return nil, errors.New("could not parse latitude in TML: " + err.Error())
		}
		lon, err := Parse45Digit(segments[i+1])
		if err != nil {
			return nil, errors.New("could not parse longitude in TML: " + err.Error())
		}
		p, err := geom.NewPoint(geom.XY).SetCoords([]float64{-float64(lon), float64(lat)})
		if err != nil {
			return nil, errors.New("could not create point in TML: " + err.Error())
		}
		err = points.Push(p)
		if err != nil {
			return nil, errors.New("could not push point to multipoint in TML: " + err.Error())
		}
	}

	tml := TML{
		Original:    original,
		Time:        time,
		Direction:   direction,
		Speed:       speed,
		SpeedString: speedString,
		Locations:   points,
	}

	return &tml, nil
}
