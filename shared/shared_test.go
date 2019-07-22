package shared

import "testing"

func TestCleanFileName(t *testing.T) {

	testsString := "+testN?a%m#e~"
	cleanedString := CleanFileName(testsString)
	if cleanedString != "testName" {
		t.Errorf("leanedString was incorrect, got: %s, want: %s.", cleanedString, "testName")
	}
}

func TestGetSetting(t *testing.T) {

	setting := GetSetting()

	if len(setting.Engine) == 0 {
		t.Errorf("setting.Engine is empty!")
	}

	if len(setting.Term) == 0 {
		t.Errorf("setting.Term is empty!")
	}

	if len(setting.FolderName) == 0 {
		t.Errorf("setting.FolderName is empty!")
	}

	if setting.Limit == 0 {
		t.Errorf("setting.Limit is zero!")
	}

	if len(setting.ImageType) == 0 {
		t.Errorf("setting.ImageType is empty!")
	}

	if len(setting.FolderPath) == 0 {
		t.Errorf("setting.FolderPath is empty!")
	}
}
