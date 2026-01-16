package awips

import "testing"

func TestFindAWIPS(t *testing.T) {
	expected := "ZFPLWX"

	// Valid AWIPS
	text := "ZFPLWX\n"
	if h := FindAWIPS(text); h != expected {
		t.Errorf("expected %s, got %s", expected, h)
	}

	expected = "SEL0"

	// Valid AWIPS with spaces
	text = "SEL0  \n"
	if h := FindAWIPS(text); h != expected {
		t.Errorf("expected %s, got %s", expected, h)
	}

	// Valid AWIPS multiline
	text = "\nSEL0\n\n"
	if h := FindAWIPS(text); h != expected {
		t.Errorf("expected %s, got %s", expected, h)
	}

	// For invalid or not found AWIPS, we expect an empty string
	expected = ""

	// Inalid AWIPS, invalid leading space
	text = " SEL0"
	if h := FindAWIPS(text); h != expected {
		t.Errorf("expected %s, got %s", expected, h)
	}

	// Inalid AWIPS, invalid trailing spaces
	text = "SEL0   "
	if h := FindAWIPS(text); h != expected {
		t.Errorf("expected %s, got %s", expected, h)
	}

	// Invalid AWIPS, invalid characters
	text = "sel0\n"
	if h := FindAWIPS(text); h != expected {
		t.Errorf("expected %s, got %s", expected, h)
	}

	text = "SVRTHUN\n"
	if h := FindAWIPS(text); h != expected {
		t.Errorf("expected %s, got %s", expected, h)
	}
}

func TestHasAWIPS(t *testing.T) {
	expected := true

	// Full valid AWIPS
	if b := HasAWIPS("ZFPLWX\n"); b != expected {
		t.Errorf("expected %t, got %t", expected, b)
	}

	// Valid partial AWIPS
	if b := HasAWIPS("SEL0\n"); b != expected {
		t.Errorf("expected %t, got %t", expected, b)
	}

	// Valid partial AWIPS with spaces
	if b := HasAWIPS("SEL0  \n"); b != expected {
		t.Errorf("expected %t, got %t", expected, b)
	}

	expected = false

	// Invalid AWIPS
	if b := HasAWIPS("sel  \n"); b != expected {
		t.Errorf("expected %t, got %t", expected, b)
	}
}

func TestParseAWIPS(t *testing.T) {
	expected := AWIPS{
		Original: "ZFPLWX",
		Product:  "ZFP",
		NWSLI:    "LWX",
	}

	// Valid AWIPS
	text := "ZFPLWX\n"
	if awips, err := ParseAWIPS(text); err != nil || awips != expected {
		t.Errorf("expected %v, got %v, error: %v", expected, awips, err)
	}

	// Valid AWIPS with spaces
	text = "SEL0  \n"
	expected = AWIPS{
		Original: "SEL0",
		Product:  "SEL",
		NWSLI:    "0",
	}
	if awips, err := ParseAWIPS(text); err != nil || awips != expected {
		t.Errorf("expected %v, got %v, error: %v", expected, awips, err)
	}

	// Invalid AWIPS
	text = "sel  \n"
	if _, err := ParseAWIPS(text); err == nil {
		t.Error("expected an error for invalid AWIPS")
	}
}
