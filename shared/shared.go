package shared

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

type Setting struct {
	Engine     string
	Term       string
	FolderName string
	Limit      int
	ImageType  string
	FolderPath string
}

func GetSetting() (Setting, error) {
	engine := flag.String("engine", "yandex", "search engine: google, bing, or yandex")
	folderName := flag.String("folder", "img", "destination folder name")
	term := flag.String("term", "apple", "search term")
	limit := flag.Int("limit", 75, "maximum number of images to download")
	imageType := flag.String("type", "*", "image file type filter (e.g., jpeg, png, or *)")
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return Setting{}, err
	}

	folderPath := filepath.Join(dir, *folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
			return Setting{}, err
		}
	}

	var setting Setting
	setting.Term = strings.Replace(*term, " ", "+", -1)
	setting.FolderName = *folderName
	setting.FolderPath = folderPath
	setting.Limit = *limit
	setting.ImageType = *imageType
	setting.Engine = *engine

	return setting, nil
}

func CleanFileName(fileName string) string {
	symbols := [10]string{"*", "?", "%", "\\", "/", " ", "+", "#", "@", "~"}
	for _, symbol := range symbols {
		fileName = strings.Replace(fileName, symbol, "", -1)
	}
	return fileName
}

func PrintError(err error) {
	println("<==============ERROR======================>")
	println(err.Error())
	println("<==============ERROR======================>")
}
