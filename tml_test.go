package awips

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindTML(t *testing.T) {
	// Example from NWS Directive 10-1701
	text := `TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318
	`
	tml := FindTML(text)
	assert.NotEqual(t, "", tml)
	assert.Equal(t, "TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318", tml)

	// Multiple points
	text = `TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318 3490 10420
	`
	tml = FindTML(text)
	assert.NotEqual(t, "", tml)
	assert.Equal(t, "TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318 3490 10420", tml)

	// Multiple lines
	text = `TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318 
    3490 10420`
	tml = FindTML(text)
	assert.NotEqual(t, "", tml)
	assert.Equal(t, `TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318 
    3490 10420`, tml)

	// Lowercase prefix (should not match)
	text = `time...mot...loc 0128Z 004DEG 9KT 3480 10318`
	tml = FindTML(text)
	assert.Equal(t, "", tml)
}

func TestParseTML(t *testing.T) {
	issued := time.Date(2025, 8, 21, 0, 0, 0, 0, time.UTC)

	// Valid TML
	text := `TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318`
	tml, err := ParseTML(text, issued)
	assert.NoError(t, err)
	assert.NotNil(t, tml)
	assert.Equal(t, "TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318", tml.Original)
	assert.Equal(t, time.Date(2025, 8, 21, 1, 28, 0, 0, time.UTC), tml.Time)
	assert.Equal(t, 4, tml.Direction)
	assert.Equal(t, 9, tml.Speed)
	assert.Equal(t, "9KT", tml.SpeedString)
	assert.Equal(t, 1, tml.Locations.NumCoords())
	assert.Equal(t, -103.18, tml.Locations.Coord(0).X())
	assert.Equal(t, 34.8, tml.Locations.Coord(0).Y())

	// Multiple locations
	text = `TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318 3490 10420`
	tml, err = ParseTML(text, issued)
	assert.NoError(t, err)
	assert.NotNil(t, tml)
	assert.Equal(t, "TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318 3490 10420", tml.Original)
	assert.Equal(t, 2, tml.Locations.NumCoords())
	assert.Equal(t, -104.2, tml.Locations.Coord(1).X())
	assert.Equal(t, 34.9, tml.Locations.Coord(1).Y())

	// Multi Line TML
	text = `TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318
    3490 10420`
	tml, err = ParseTML(text, issued)
	assert.NoError(t, err)
	assert.NotNil(t, tml)
	assert.Equal(t, "TIME...MOT...LOC 0128Z 004DEG 9KT 3480 10318 3490 10420", tml.Original)
	assert.Equal(t, 2, tml.Locations.NumCoords())
	assert.Equal(t, -103.18, tml.Locations.Coord(0).X())
	assert.Equal(t, 34.8, tml.Locations.Coord(0).Y())
	assert.Equal(t, 2, tml.Locations.NumCoords())
	assert.Equal(t, -104.2, tml.Locations.Coord(1).X())
	assert.Equal(t, 34.9, tml.Locations.Coord(1).Y())

	// Motion less than 1
	text = `TIME...MOT...LOC 0128Z 004DEG 0KT 3480 10318`
	tml, err = ParseTML(text, issued)
	assert.NoError(t, err)
	assert.NotNil(t, tml)
	assert.Equal(t, "TIME...MOT...LOC 0128Z 004DEG 0KT 3480 10318", tml.Original)
	assert.Equal(t, 0, tml.Speed) // Speed should be 0

	// Some products have the motion in a unit other than knots. Test for this.
	text = `TIME...MOT...LOC 0128Z 004DEG 9MPH 3480 10318`
	tml, err = ParseTML(text, issued)
	assert.NoError(t, err)
	assert.NotNil(t, tml)
	assert.Equal(t, "TIME...MOT...LOC 0128Z 004DEG 9MPH 3480 10318", tml.Original)
	assert.Equal(t, 9, tml.Speed)            // Speed should still be parsed as knots
	assert.Equal(t, "9MPH", tml.SpeedString) // Speed string should reflect
}

// text := `TIME...MOT...LOC 1300Z 090DEG 20KT 3881 10015`
