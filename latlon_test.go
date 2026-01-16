package awips

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twpayne/go-geom"
)

// 4-5 digit values are usually used by WFOs for local products.
func Test4Digit(t *testing.T) {
	// 48.96 N, valid
	f, err := Parse45Digit("4896")
	assert.NoError(t, err)
	assert.Equal(t, 48.96, f)

	// 80.00 W, valid
	f, err = Parse45Digit("8000")
	assert.NoError(t, err)
	assert.Equal(t, 80.00, f)

	// 179.99 W, valid
	f, err = Parse45Digit("17999")
	assert.NoError(t, err)
	assert.Equal(t, 179.99, f)

	// 180.0 W, valid
	f, err = Parse45Digit("18000")
	assert.NoError(t, err)
	assert.Equal(t, 180.0, f)

	// 180.01 W, valid
	f, err = Parse45Digit("18001")
	assert.NoError(t, err)
	assert.Equal(t, 180.01, f)

	// Invalid length
	_, err = Parse45Digit("48")
	assert.Error(t, err)
	_, err = Parse45Digit("123456")
	assert.Error(t, err)

	// Invalid characters
	_, err = Parse45Digit("48a6")
	assert.Error(t, err)
}

// 8-digits are usually used for national products.
func Test8Digit(t *testing.T) {
	// 48.96 N, 80.00 W, valid
	c, err := ParseCoord8("48968000")
	assert.NoError(t, err)
	assert.Equal(t, geom.Coord{80.0, 48.96}, c)

	// 48.96 N, 99.99 E, valid
	c, err = ParseCoord8("48969999")
	assert.NoError(t, err)
	assert.Equal(t, geom.Coord{99.99, 48.96}, c)

	// 48.96 N, 100.00 E, valid
	c, err = ParseCoord8("48960000")
	assert.NoError(t, err)
	assert.Equal(t, geom.Coord{100.0, 48.96}, c)

	// 48.96 N, 100.01 E, valid
	c, err = ParseCoord8("48960001")
	assert.NoError(t, err)
	assert.Equal(t, geom.Coord{100.01, 48.96}, c)

	// Invalid length
	_, err = ParseCoord8("4896")
	assert.Error(t, err)
	_, err = ParseCoord8("4896800")
	assert.Error(t, err)

	// Invalid characters
	_, err = ParseCoord8("4896a000")
	assert.Error(t, err)
}

func TestParseLatLon(t *testing.T) {
	// Valid 4 digit LAT...LON text
	text := `LAT...LON 3468 9976 3468 9967 3454 9967 3451 9983
      3456 9989`
	latlon, err := ParseLatLon(text)
	assert.NoError(t, err)
	assert.Equal(t, text, latlon.Original)
	assert.Equal(t, 5, len(latlon.Coords))
	assert.Equal(t, []geom.Coord{{99.76, 34.68}, {99.67, 34.68}, {99.67, 34.54}, {99.83, 34.51}, {99.89, 34.56}}, latlon.Coords)

	// Valid 4-5 digit LAT...LON text
	text = `LAT...LON 2069 15600 2064 15607 2059 15642 2064 15646
      2079 15647 2081 15663 2089 15669 2096 15669
      2101 15666 2104 15660 2090 15648 2095 15633
      2094 15625 2082 15611 2080 15601 2070 15600`
	latlon, err = ParseLatLon(text)
	assert.NoError(t, err)
	assert.Equal(t, text, latlon.Original)
	assert.Equal(t, 16, len(latlon.Coords))
	assert.Equal(t, []geom.Coord{{156.00, 20.69}, {156.07, 20.64}, {156.42, 20.59}, {156.46, 20.64},
		{156.47, 20.79}, {156.63, 20.81}, {156.69, 20.89}, {156.69, 20.96},
		{156.66, 21.01}, {156.60, 21.04}, {156.48, 20.9}, {156.33, 20.95},
		{156.25, 20.94}, {156.11, 20.82}, {156.01, 20.8}, {156.0, 20.7}}, latlon.Coords)

	// Valid 8 digit LAT...LON text
	text = `LAT...LON   36740503 40900506 40910219 
	  36740234 36740503`
	latlon, err = ParseLatLon(text)
	assert.NoError(t, err)
	assert.Equal(t, text, latlon.Original)
	assert.Equal(t, 5, len(latlon.Coords))
	assert.Equal(t, []geom.Coord{{105.03, 36.74}, {105.06, 40.90}, {102.19, 40.91}, {102.34, 36.74}, {105.03, 36.74}}, latlon.Coords)

	// Invalid length LAT...LON text
	text = `LAT...LON 1234 5678 91011`
	latlon, err = ParseLatLon(text)
	assert.Error(t, err)

	// Invalid number of characters LAT...LON text
	text = `LAT...LON 1234 5678 91011 123456`
	latlon, err = ParseLatLon(text)
	assert.Error(t, err)
}

func TestLonWestBias(t *testing.T) {
	// Test we get the raw information without the west-bias applied.
	text := `LAT...LON 2000 17900 2000 18100`
	latlon, err := ParseLatLon(text)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(latlon.Coords))
	assert.Equal(t, []geom.Coord{{179.00, 20.00}, {181.00, 20.00}}, latlon.Coords)

	// Test the west-bias is applied correctly.
	west := LonWestBias(latlon.Coords)
	assert.Equal(t, 2, len(latlon.Coords))
	assert.Equal(t, []geom.Coord{{179.00, 20.00}, {-179.00, 20.00}}, latlon.Coords)
	assert.Equal(t, 2, len(west))
	assert.Equal(t, []geom.Coord{{179.00, 20.00}, {-179.00, 20.00}}, west)
}

func TestCloseRing(t *testing.T) {
	// Test we get the raw information without the closed ring.
	text := `LAT...LON 2000 17900 2000 18100`
	latlon, err := ParseLatLon(text)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(latlon.Coords))
	assert.Equal(t, []geom.Coord{{179.00, 20.00}, {181.00, 20.00}}, latlon.Coords)

	// Test we get the closed ring
	latlon, err = latlon.ClosedRing()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(latlon.Coords))
	assert.Equal(t, []geom.Coord{{179.00, 20.00}, {181.00, 20.00}, {179.00, 20.00}}, latlon.Coords)

	// Test empty coords
	latlon = &LatLon{}
	latlon, err = latlon.ClosedRing()
	assert.Error(t, err)
}
