package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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

func main() {

	imageFolder := flag.String("folder", "img", "a string")
	searchTerm := flag.String("searchTerm", "image", "a string")
	flag.Parse()

	println("Folder:" + *imageFolder)
	println("Search term:" + *searchTerm)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	folderPath := dir + "\\" + *imageFolder
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, os.ModePerm)
	}

	//Preparing term for the search
	term := strings.Replace(*searchTerm, " ", "+", -1)

	url := "https://www.google.com/search?q=" + term + "&oq=" + term + "&biw=1536&bih=723&tbm=isch&sa=1&ei=6qqGXM_oDenYjwSw1b-oAw"
	//url = 'https://www.google.com/search?q=' + keywordem          + &espv=2&biw=1366&bih=667&site=webhp&source=lnms&tbm=isch&sa=X&ei=XosDVaCXD8TasATItgE&ved=0CAcQ_AUoAg'
	//url = 'https://www.google.com/search?q=' + quote(search_term) + &espv=2&biw=1366&bih=667&site=webhp&source=lnms&tbm=isch' + params + '&sa=X&ei=XosDVaCXD8TasATItgE&ved=0CAcQ_AUoAg'
	//url = 'https://www.google.com/search?q=' + quote(search_term) + &espv=2&biw=1366&bih=667&site=webhp&source=lnms&tbm=isch' + params + '&sa=X&ei=XosDVaCXD8TasATItgE&ved=0CAcQ_AUoAg'

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	imageSource := string(body)

	r := regexp.MustCompile("<div class=\"rg_meta notranslate\">([\\s\\S]*?)</div>")
	matches := r.FindAllStringSubmatch(imageSource, -1)

	println("GOT ITEMS: ", len(matches))

	for _, each := range matches {

		img := GetImageItemFromJson(each[1])
		fullName := GetFileFullName(img, folderPath)

		if err := DownloadImage(fullName, img.Ou); err != nil {
			//panic(err)
			println("<==============ERROR======================>")
			println(err.Error())
			println(each[1])
			println(img.Ou)
			println("<==============ERROR======================>")
		}
	}
}

func GetFileFullName(img ImageItem, folderPath string) string {
	url := img.Ou
	fileName := ""
	if len(img.Ity) != 0 {
		fileName = img.ID[0:len(img.ID)-1] + "_" + url[strings.LastIndex(img.Ou, "/")+1:strings.Index(img.Ou, "."+img.Ity)] + "." + img.Ity
	} else {
		fileName = img.ID[0:len(img.ID)-1] + ".jpeg"
	}

	fileName = folderPath + "\\" + CleanFileName(fileName)
	//println("<===================================================>")
	// println("-----------------------------------------------------")
	// println("Source:", img.Ou)
	println(fileName)
	// println(img.Ou)
	//println("<===================================================>")

	return fileName
}

func CleanFileName(fileName string) string {
	symbols := [6]string{"*", "?", "%", "\\", "/"}
	for _, symbol := range symbols {
		fileName = strings.Replace(fileName, symbol, "", -1)
	}
	return fileName
}

func GetImageItemFromJson(jsonString string) ImageItem {

	// fmt.Println("-----------------------------------------")
	// fmt.Println(json)
	// fmt.Println("-----------------------------------------")

	img := ImageItem{}

	err := json.Unmarshal([]byte(jsonString), &img)
	if err != nil {
		panic(err)
	}

	return img
}

func DownloadImage(filePath string, url string) error {

	// Get the data
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("authority", "cdn-images-1.medium.com")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
