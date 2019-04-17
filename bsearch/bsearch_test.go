package bsearch

import "testing"

func TestCleanFileName(t *testing.T) {
	jsonImage := "{\"cid\":\"AB5rIo1B\",\"purl\":\"http://www.michiganapples.com/\",\"murl\":\"http://www.michiganapples.com/Portals/0/New%20Fall.jpg\",\"turl\":\"https://tse3.mm.bing.net/th?id=OIP.AB5rIo1B5wT8qLNu3dg4KAHaJ4&amp;pid=15.1\",\"md5\":\"001e6b228d41e704fca8b36eddd83828\",\"shkey\":\"\",\"t\":\"Michigan Apple Committee\",\"mid\":\"E52BA92E1691011E3111A4106B54E5AB5AC3898B\",\"desc\":\"\"}"
	img, err := GetImageItemFromJson(jsonImage)

	if err != nil {
		t.Errorf("Incorrect json")
	}

	if len(img.Murl) == 0 {
		t.Errorf("Image url is empty")
	}
}

func TestCleanFileNameEmptyJson(t *testing.T) {
	jsonImage := " "
	img, err := GetImageItemFromJson(jsonImage)

	if err == nil && len(img.Murl) == 0 {
		t.Errorf("GetImageItemFromJson fails")
	}
}
