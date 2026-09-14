package extensions

import (
	"strconv"
	"strings"
)

func validSemver(s string) bool {
	if len(s) == 0 || len(s) > MaxVersionBytes {
		return false
	}
	mainAndBuild := strings.SplitN(s, "+", 2)
	if len(mainAndBuild) == 2 && !validIdentifierList(mainAndBuild[1], false) {
		return false
	}
	coreAndPre := strings.SplitN(mainAndBuild[0], "-", 2)
	if len(coreAndPre) == 2 && !validIdentifierList(coreAndPre[1], true) {
		return false
	}
	parts := strings.Split(coreAndPre[0], ".")
	if len(parts) != 3 {
		return false
	}
	for _, part := range parts {
		if !validNumericIdentifier(part) {
			return false
		}
	}
	return true
}

func validNumericIdentifier(s string) bool {
	if s == "" || len(s) > 1 && s[0] == '0' {
		return false
	}
	for i := range len(s) {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	n, err := strconv.ParseUint(s, 10, 31)
	return err == nil && n <= maxVersionNumber
}

func validIdentifierList(s string, rejectNumericLeadingZero bool) bool {
	if s == "" {
		return false
	}
	for _, part := range strings.Split(s, ".") {
		if part == "" {
			return false
		}
		numeric := true
		for i := range len(part) {
			c := part[i]
			if !isASCIILetter(c) && (c < '0' || c > '9') && c != '-' {
				return false
			}
			if c < '0' || c > '9' {
				numeric = false
			}
		}
		if rejectNumericLeadingZero && numeric && len(part) > 1 && part[0] == '0' {
			return false
		}
	}
	return true
}

func isASCIILetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
