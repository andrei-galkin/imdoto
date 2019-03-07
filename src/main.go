package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type logWriter struct{}

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

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	folderPath := dir + "\\img"
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, os.ModePerm)
	}

	url := "https://www.google.com/search?q=c%23+image&rlz=1C1GCEU_enCA826CA826&source=lnms&tbm=isch&sa=X&ved=0ahUKEwjfgeT_pu7gAhVp74MKHZZHB9wQ_AUIDigB&biw=1536&bih=723&dpr=1.25#imgrc=40_i-tVgNtfjAM:"

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

	fmt.Println("\n", len(imageSource))

	r := regexp.MustCompile("<div class=\"rg_meta notranslate\">([\\s\\S]*?)</div>")

	matches := r.FindAllStringSubmatch(imageSource, -1)
	fmt.Println("GOT ITEMS: ", len(matches))

	for _, each := range matches {

		jsonString := each[1]
		// fmt.Println("-----------------------------------------")
		// fmt.Println(jsonString)
		// fmt.Println("-----------------------------------------")
		img := ImageItem{}

		err := json.Unmarshal([]byte(jsonString), &img)
		if err != nil {
			panic(err)
		}

		url := img.Ou

		fileName := ""

		if len(img.Ity) != 0 {
			fileName = img.ID[0:len(img.ID)-1] + "_" + url[strings.LastIndex(img.Ou, "/")+1:strings.Index(img.Ou, "."+img.Ity)] + "." + img.Ity
		} else {
			fileName = img.ID[0:len(img.ID)-1] + ".jpeg"
		}

		fullName := folderPath + "\\" + fileName
		// println("<===================================================>")
		// println("-----------------------------------------------------")
		// println("Source:", img.Ou)
		// println(fileName)
		// println(img.Ou)
		// println("<===================================================>")

		if err := DownloadImage(fullName, img.Ou); err != nil {
			//panic(err)
			println("<==============ERROR======================>")
			println(err.Error())
			println(jsonString)
			println(img.Ou)
			println("<==============ERROR======================>")
		}
	}
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
