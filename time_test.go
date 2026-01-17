package awips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTimezones(t *testing.T) {

	tz := GetTimezone("GMT")
	assert.NotNil(t, tz)
	tz = GetTimezone("gmt")
	assert.NotNil(t, tz)

	tz = GetTimezone("UTC")
	assert.NotNil(t, tz)
	tz = GetTimezone("utc")
	assert.NotNil(t, tz)

	tz = GetTimezone("EST")
	assert.NotNil(t, tz)
	tz = GetTimezone("est")
	assert.NotNil(t, tz)

	tz = GetTimezone("EDT")
	assert.NotNil(t, tz)
	tz = GetTimezone("edt")
	assert.NotNil(t, tz)

	// Special Cases
	tz = GetTimezone("ChST")
	assert.NotNil(t, tz)
	tz = GetTimezone("CHST")
	assert.NotNil(t, tz)
	tz = GetTimezone("chst")
	assert.NotNil(t, tz)
}
