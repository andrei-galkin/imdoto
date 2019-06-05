package bingsearch

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

func TestRestoreQuotes(t *testing.T) {
	jsonImage := "{&quot;cid&quot;:&quot;AB5rIo1B&quot;,&quot;purl&quot;:&quot;http://www.michiganapples.com/&quot;,&quot;murl&quot;:&quot;http://www.michiganapples.com/Portals/0/New%20Fall.jpg&quot;,&quot;turl&quot;:&quot;https://tse3.mm.bing.net/th?id=OIP.AB5rIo1B5wT8qLNu3dg4KAHaJ4&amp;pid=15.1&quot;,&quot;md5&quot;:&quot;001e6b228d41e704fca8b36eddd83828&quot;,&quot;shkey&quot;:&quot;&quot;,&quot;t&quot;:&quot;Michigan Apple Committee&quot;,&quot;mid&quot;:&quot;E52BA92E1691011E3111A4106B54E5AB5AC3898B&quot;,&quot;desc&quot;:&quot;&quot;}"
	jsonImage = RestoreQuotes(jsonImage)

	img, err := GetImageItemFromJson(jsonImage)

	if err != nil {
		t.Errorf("Incorrect json")
	}

	if len(img.Murl) == 0 {
		t.Errorf("Image url is empty")
	}
}

func TestGetFileFullName(t *testing.T) {
	folderPath := `c:`
	expectedFullName := `c:\AB5rIo1B_New20Fall.jpg`
	jsonImage := "{&quot;cid&quot;:&quot;AB5rIo1B&quot;,&quot;purl&quot;:&quot;http://www.michiganapples.com/&quot;,&quot;murl&quot;:&quot;http://www.michiganapples.com/Portals/0/New%20Fall.jpg&quot;,&quot;turl&quot;:&quot;https://tse3.mm.bing.net/th?id=OIP.AB5rIo1B5wT8qLNu3dg4KAHaJ4&amp;pid=15.1&quot;,&quot;md5&quot;:&quot;001e6b228d41e704fca8b36eddd83828&quot;,&quot;shkey&quot;:&quot;&quot;,&quot;t&quot;:&quot;Michigan Apple Committee&quot;,&quot;mid&quot;:&quot;E52BA92E1691011E3111A4106B54E5AB5AC3898B&quot;,&quot;desc&quot;:&quot;&quot;}"
	jsonImage = RestoreQuotes(jsonImage)

	img, err := GetImageItemFromJson(jsonImage)
	fullName := GetFileFullName(img, folderPath)
	if err != nil {
		t.Errorf("Incorrect json")
	}

	if len(img.Murl) == 0 {
		t.Errorf("Image url is empty")
	}

	if fullName != expectedFullName {
		t.Errorf("%s fullName is incorrect, expected %s", fullName, expectedFullName)
	}
}
