package njalla

import (
	"fmt"
	"testing"
)

func TestParseImportIDExpected(t *testing.T) {
	expectedDomain := "testing.com"
	expectedID := "1234"
	input := fmt.Sprintf("%s:%s", expectedDomain, expectedID)

	resultDomain, resultID, err := parseImportID(input)
	if err != nil {
		t.Fatalf("%q", err)
	}

	if resultDomain != expectedDomain {
		t.Fatalf(
			"Result domain %s doesn't match expected domain %s",
			resultDomain, expectedDomain,
		)
	}

	if resultID != expectedID {
		t.Fatalf(
			"Result ID %s doesn't match expected ID %s",
			resultID, expectedID,
		)
	}
}

func TestParseImportIDMissingIDWithColon(t *testing.T) {
	expectedDomain := "testing.com"
	input := fmt.Sprintf("%s:", expectedDomain)

	_, _, err := parseImportID(input)
	if err == nil {
		t.Fatal("Unexpected success")
	}
}

func TestParseImportIDMissingIDWithoutColon(t *testing.T) {
	expectedDomain := "testing.com"
	input := fmt.Sprintf("%s", expectedDomain)

	_, _, err := parseImportID(input)
	if err == nil {
		t.Fatal("Unexpected success")
	}
}
