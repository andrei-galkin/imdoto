package imdoto

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

func GetSetting() Setting {
	engine := flag.String("engine", "bing", "a string")
	folderName := flag.String("folder", "img", "a string")
	term := flag.String("term", "apple", "a string")
	limit := flag.Int("limit", 700, "a int")
	imageType := flag.String("type", "*", "a string")
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	folderPath := dir + `\` + *folderName
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, os.ModePerm)
	}

	var setting Setting
	setting.Term = strings.Replace(*term, " ", "+", -1)
	setting.FolderName = *folderName
	setting.FolderPath = folderPath
	setting.Limit = *limit
	setting.ImageType = *imageType
	setting.Engine = *engine

	return setting
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
