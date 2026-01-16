package awips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidUGC(t *testing.T) {
	// Modified example from NWS Directive 10-1702
	text := `WYZ001>020-021-022>030-035-081700-`

	ugc, err := ParseUGC(text)
	assert.NoError(t, err)
	// Should just by WY (Wyoming)
	assert.Equal(t, 1, len(ugc.States))
	// ID should be WY
	assert.Equal(t, "WY", ugc.States[0].ID)
	// Type should be Z (zone)
	assert.Equal(t, "Z", ugc.States[0].Type)
	// Should have 31 zones
	assert.Equal(t, 31, len(ugc.States[0].Areas))
	//
	assert.Equal(t, []string{"001", "002", "003", "004", "005", "006", "007", "008", "009", "010", "011", "012",
		"013", "014", "015", "016", "017", "018", "019", "020", "021", "022", "023", "024",
		"025", "026", "027", "028", "029", "030", "035"}, ugc.States[0].Areas)

	// Multiple states
	text = `WYC001>020-021-022>030-035-FLC020-202200-`
	ugc, err = ParseUGC(text)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ugc.States))
	// Second state should be FL (Florida)
	assert.Equal(t, "FL", ugc.States[1].ID)
	// Type should be C (county)
	assert.Equal(t, "C", ugc.States[1].Type)
	// Should just be 1 county
	assert.Equal(t, 1, len(ugc.States[1].Areas))
	assert.Equal(t, []string{"020"}, ugc.States[1].Areas)
}

func TestInvalidUGC(t *testing.T) {
	// Missing datetime
	text := `WYZ001-`
	ugc, err := ParseUGC(text)
	assert.NoError(t, err)
	assert.Nil(t, ugc)

	// Invalid County/Zone format code
	text = `CAF001-081700-`
	ugc, err = ParseUGC(text)
	assert.NoError(t, err)
	assert.Nil(t, ugc)

	// Invalid datetime format
	text = `WYZ001-0817-`
	ugc, err = ParseUGC(text)
	assert.NoError(t, err)
	assert.Nil(t, ugc)

	// Invalid datetime date
	text = `WYZ001-001700-`
	_, err = ParseUGC(text)
	assert.Error(t, err)

	// Invalid datetime time
	text = `FLC020-012500-`
	_, err = ParseUGC(text)
	assert.Error(t, err)
}
