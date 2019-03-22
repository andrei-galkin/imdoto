package imdoto

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

type ImageItem struct {
	ID  string `json:"id"`
	Isu string `json:"isu"`
	Itg int    `json:"itg"`
	Ity string `json:"ity"`
	Oh  int    `json:"oh"`
	Ou  string `json:"ou"`
	Ow  int    `json:"ow"`
	Pt  string `json:"pt"`
	Rh  string `json:"rh"`
	Rid string `json:"rid"`
	Rt  int    `json:"rt"`
	Ru  string `json:"ru"`
	S   string `json:"s"`
	Sc  int    `json:"sc"`
	St  string `json:"st"`
	Th  int    `json:"th"`
	Tu  string `json:"tu"`
	Tw  int    `json:"tw"`
}

type LoadOption struct {
	Term       string
	FolderName string
	Limit      int
	ImageType  string
	FolderPath string
}

func GetLoadOption() LoadOption {
	folderName := flag.String("folder", "img", "a string")
	term := flag.String("term", "apple fruit", "a string")
	limit := flag.Int("limit", 12, "a int")
	imageType := flag.String("type", "*", "a string")
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	folderPath := dir + "\\" + *folderName
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, os.ModePerm)
	}

	var option LoadOption
	option.Term = strings.Replace(*term, " ", "+", -1)
	option.FolderName = *folderName
	option.FolderPath = folderPath
	option.Limit = *limit
	option.ImageType = *imageType

	return option
}

func CleanFileName(fileName string) string {
	symbols := [6]string{"*", "?", "%", "\\", "/"}
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
