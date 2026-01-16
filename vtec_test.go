package awips

import (
	"testing"
	"time"
)

func TestVTECParse(t *testing.T) {
	vtecs, err := ParseVTEC("/O.NEW.KRAH.SV.W.0175.250626T2336Z-250627T0015Z/")
	if len(err) > 0 {
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}

	// Expect one VTEC
	if len(vtecs) != 1 {
		t.Errorf("expected 1 VTEC, got %d", len(vtecs))
	}

	vtec := vtecs[0]

	// Test VTEC fields
	if vtec.Class != "O" {
		t.Errorf("expected class 'O', got '%s'", vtec.Class)
	}
	if vtec.Action != "NEW" {
		t.Errorf("expected action 'NEW', got '%s'", vtec.Action)
	}
	if vtec.WFO != "KRAH" {
		t.Errorf("expected WFO 'KRAH', got '%s'", vtec.WFO)
	}
	if vtec.Phenomena != "SV" {
		t.Errorf("expected phenomena 'SV', got '%s'", vtec.Phenomena)
	}
	if vtec.Significance != "W" {
		t.Errorf("expected significance 'W', got '%s'", vtec.Significance)
	}
	if vtec.EventNumber != 175 {
		t.Errorf("expected event number 175, got %d", vtec.EventNumber)
	}
	expectedStart := time.Date(2025, time.June, 26, 23, 36, 0, 0, time.UTC)
	if !vtec.Start.Equal(expectedStart) {
		t.Errorf("expected start time '%s', got '%s'", expectedStart, vtec.Start)
	}
	expectedEnd := time.Date(2025, time.June, 27, 0, 15, 0, 0, time.UTC)
	if !vtec.End.Equal(expectedEnd) {
		t.Errorf("expected end time '%s', got '%s'", expectedEnd, vtec.End)
	}
	if vtec.Original != "O.NEW.KRAH.SV.W.0175.250626T2336Z-250627T0015Z" {
		t.Errorf("expected original 'O.NEW.KRAH.SV.W.0175.250626T2336Z-250627T0015Z', got '%s'", vtec.Original)
	}
}

func TestVTECParseMultiple(t *testing.T) {
	vtecs, err := ParseVTEC(`
	/O.CAN.KDMX.TO.W.0045.000000T0000Z-240521T2145Z/
	/O.CON.KDMX.TO.W.0045.000000T0000Z-240521T2145Z/
	`)
	if len(err) > 0 {
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}

	if len(vtecs) != 2 {
		t.Errorf("expected 2 VTECs, got %d", len(vtecs))
	}

	if vtecs[0].Action != "CAN" {
		t.Errorf("expected first VTEC action 'CAN', got '%s'", vtecs[0].Action)
	}
	if vtecs[1].Action != "CON" {
		t.Errorf("expected second VTEC action 'CON', got '%s'", vtecs[1].Action)
	}
}

func TestVTECNilDates(t *testing.T) {
	vtecs, err := ParseVTEC(`
	/O.EXT.KJKL.HT.Y.0001.000000T0000Z-250628T0000Z/
	/O.NEW.KJKL.FL.W.0001.250628T0000Z-000000T0000Z/
	`)
	if len(err) > 0 {
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}

	// Expect one VTEC
	if len(vtecs) != 2 {
		t.Errorf("expected 2 VTECs, got %d", len(vtecs))
	}

	// Test VTEC fields
	if vtecs[0].Start != nil {
		t.Errorf("expected start time to be nil, got '%s'", vtecs[0].Start)
	}
	if vtecs[1].End != nil {
		t.Errorf("expected end time to be nil, got '%s'", vtecs[1].End)
	}
}

func TestVTECParseInvalid(t *testing.T) {
	// Completely wrong
	_, err := ParseVTEC("/S.UPD.KJKL.JJ.K.0000.251332T2461Z-251332T2461Z/")
	if len(err) != 1 {
		t.Errorf("expected 1 error, got %d", len(err))
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}

	// Check with valid class
	_, err = ParseVTEC("/O.UPD.KJKL.JJ.K.0000.251332T2461Z-251332T2461Z/")
	if len(err) != 1 {
		t.Errorf("expected 1 error, got %d", len(err))
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}

	// Check with valid action
	_, err = ParseVTEC("/O.CON.KJKL.SV.K.0000.251332T2461Z-251332T2461Z/")
	if len(err) != 1 {
		t.Errorf("expected 1 error, got %d", len(err))
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}

	// Check with valid phenomena
	_, err = ParseVTEC("/O.CON.KJKL.SV.K.0000.251332T2461Z-251332T2461Z/")
	if len(err) != 1 {
		t.Errorf("expected 1 error, got %d", len(err))
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}

	// Check with valid significance
	_, err = ParseVTEC("/O.CON.KJKL.SV.K.0000.251332T2461Z-251332T2461Z/")
	if len(err) != 1 {
		t.Errorf("expected 1 error, got %d", len(err))
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}

	// Check with valid event number
	_, err = ParseVTEC("/O.CON.KJKL.SV.W.0001.251332T2461Z-251332T2461Z/")
	if len(err) != 1 {
		t.Errorf("expected 1 error, got %d", len(err))
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}

	// Check with valid start time
	_, err = ParseVTEC("/O.CON.KJKL.SV.W.0001.000000T0000Z-251332T2461Z/")
	if len(err) != 1 {
		t.Errorf("expected 1 error, got %d", len(err))
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}

	// Check with all valid
	_, err = ParseVTEC("/O.CON.KJKL.SV.W.0001.000000T0000Z-250628T0000Z/")
	if len(err) != 0 {
		t.Errorf("expected 0 errors, got %d", len(err))
		for _, e := range err {
			t.Errorf("failed to parse VTEC: %v", e)
		}
	}
}

func TestVTECTitleSpecialCases(t *testing.T) {
	// Special Marine Warning
	vtec := VTEC{
		Phenomena:    "MA",
		Significance: "W",
	}

	if vtec.Title(false) != "Special Marine Warning" {
		t.Errorf("expected title 'Special Marine Warning', got '%s'", vtec.Title(false))
	}

	// Red Flag Warning
	vtec = VTEC{
		Phenomena:    "FW",
		Significance: "W",
	}

	if vtec.Title(false) != "Red Flag Warning" {
		t.Errorf("expected title 'Red Flag Warning', got '%s'", vtec.Title(false))
	}

	// Emergency
	vtec = VTEC{
		Phenomena:    "TO",
		Significance: "W",
	}

	if vtec.Title(true) != "Tornado Emergency" {
		t.Errorf("expected title 'Tornado Emergency', got '%s'", vtec.Title(true))
	}
}
