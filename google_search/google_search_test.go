package googlesearch

import "testing"

func TestGetImageItemFromJson(t *testing.T) {
	jsonImage := "{\"cb\": 3, \"cl\": 9, \"clt\": \"n\", \"cr\": 12, \"id\": \"-TtEc9M5pE7LPM:\", \"isu\": \"apple.com\", \"itg\": 0, \"ity\": \"png\", \"oh\": 302, \"ou\": \"https://www.apple.com/ac/structured-data/images/knowledge_graph_logo.png?201606271147\", \"ow\": 302, \"pt\": \"Buy Apple Pencil - Apple (CA)\", \"rh\": \"apple.com\", \"rid\": \"dQvg6kjka863yM\", \"rt\": 0, \"ru\": \"https://www.apple.com/ca/shop/product/MK0C2AM/A/apple-pencil\", \"s\": \"\", \"sc\": 1, \"st\": \"Apple\", \"th\": 225, \"tu\": \"https://encrypted-tbn0.gstatic.com/images?q\u003dtbn:ANd9GcSYfgh6R22EUz-8EUwh-e68sSEFOHWXeZHuysRW_DPkhahMgdwJ\", \"tw\": 225}"
	img, err := GetImageItemFromJson(jsonImage)

	if err != nil {
		t.Errorf("Incorrect json")
	}

	if len(img.Ou) == 0 {
		t.Errorf("Image url is empty")
	}
}

func TestGetImageItemFromEmptyJson(t *testing.T) {
	jsonImage := " "
	img, err := GetImageItemFromJson(jsonImage)

	if err == nil && len(img.Ou) == 0 {
		t.Errorf("GetImageItemFromJson fails")
	}
}

func TestGetFileFullName(t *testing.T) {
	folderPath := `c:`
	expectedFullName := `c:\-TtEc9M5pE7LPM_knowledge_graph_logo.png`
	jsonImage := "{\"cb\": 3, \"cl\": 9, \"clt\": \"n\", \"cr\": 12, \"id\": \"-TtEc9M5pE7LPM:\", \"isu\": \"apple.com\", \"itg\": 0, \"ity\": \"png\", \"oh\": 302, \"ou\": \"https://www.apple.com/ac/structured-data/images/knowledge_graph_logo.png?201606271147\", \"ow\": 302, \"pt\": \"Buy Apple Pencil - Apple (CA)\", \"rh\": \"apple.com\", \"rid\": \"dQvg6kjka863yM\", \"rt\": 0, \"ru\": \"https://www.apple.com/ca/shop/product/MK0C2AM/A/apple-pencil\", \"s\": \"\", \"sc\": 1, \"st\": \"Apple\", \"th\": 225, \"tu\": \"https://encrypted-tbn0.gstatic.com/images?q\u003dtbn:ANd9GcSYfgh6R22EUz-8EUwh-e68sSEFOHWXeZHuysRW_DPkhahMgdwJ\", \"tw\": 225}"

	img, err := GetImageItemFromJson(jsonImage)
	fullName := GetFileFullName(img, folderPath)
	if err != nil {
		t.Errorf("Incorrect json")
	}

	if len(img.Ou) == 0 {
		t.Errorf("Image url is empty")
	}

	if fullName != expectedFullName {
		t.Errorf("%s fullName is incorrect, expected %s", fullName, expectedFullName)
	}
}
