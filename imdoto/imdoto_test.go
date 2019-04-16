package imdoto

import "testing"

func TestCleanFileName(t *testing.T) {

	testsString := "+testN?a%m#e"
	cleanedString := CleanFileName(testsString)
	if cleanedString != "testName" {
		t.Errorf("leanedString was incorrect, got: %s, want: %s.", cleanedString, "testName")
	}
}
