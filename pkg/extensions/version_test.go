package extensions

import "testing"

func TestSemverValidation(t *testing.T) {
	valid := []string{
		"0.0.0", "1.2.3", "1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-0.3.7",
		"1.0.0-x.7.z.92", "1.0.0+build.1", "1.0.0-alpha+001", "2147483647.0.0",
	}
	for _, value := range valid {
		if !validSemver(value) {
			t.Fatalf("valid semantic version rejected: %q", value)
		}
	}

	invalid := []string{
		"", "1", "1.2", "1.2.3.4", "01.2.3", "1.02.3", "1.2.03",
		"2147483648.0.0", "1.0.0-", "1.0.0+", "1.0.0-alpha..1",
		"1.0.0-01", "1.0.0-alpha_1", "1.0.0+build_1", "v1.2.3",
		"1.2.3+build+again", "1.2.3-alpha-beta-!", string(make([]byte, MaxVersionBytes+1)),
	}
	for _, value := range invalid {
		if validSemver(value) {
			t.Fatalf("invalid semantic version accepted: %q", value)
		}
	}
}

func TestNumericAndIdentifierHelpers(t *testing.T) {
	if validNumericIdentifier("") || validNumericIdentifier("00") || validNumericIdentifier("a") {
		t.Fatal("invalid numeric identifier accepted")
	}
	if !validNumericIdentifier("0") || !validNumericIdentifier("123") {
		t.Fatal("valid numeric identifier rejected")
	}
	if validIdentifierList("", false) || validIdentifierList("a..b", false) || validIdentifierList("a_b", false) {
		t.Fatal("invalid identifier list accepted")
	}
	if validIdentifierList("01", true) || !validIdentifierList("01", false) || !validIdentifierList("alpha-1.beta", true) {
		t.Fatal("identifier numeric-leading-zero policy wrong")
	}
	if !isASCIILetter('a') || !isASCIILetter('Z') || isASCIILetter('1') {
		t.Fatal("ASCII letter classification wrong")
	}
}
