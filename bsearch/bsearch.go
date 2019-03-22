package bsearch

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

import im "github.com/andrei-galkin/imdoto/imdoto"

var wg sync.WaitGroup

func Download(option im.LoadOption) {
	var imageLinks []string
	imageIndex := 0

	for index := 1; index <= option.Limit; index++ {
		if imageIndex == 0 {
			imageLinks = GetImageLinks(option.Term, option.ImageType, index-1)
		}

		img, err := GetImageItemFromJson(imageLinks[imageIndex])
		if err != nil {
			im.PrintError(err)
		}

		imageIndex += 1

		wg.Add(1)
		go DownloadImage(img, option.FolderPath, index)

		//exit if there is less images then limit
		if imageIndex == len(imageLinks)-1 && len(imageLinks) != 35 {
			break
		}

		if imageIndex == 35 {
			imageIndex = 0
		}
	}

	wg.Wait()
}

func DownloadFile(filePath string, url string) error {
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

func DownloadImage(img im.ImageItem, folderPath string, index int) {
	fullName := GetFileFullName(img, folderPath)

	if err := DownloadFile(fullName, img.Ou); err != nil {
		im.PrintError(err)
	}
	indexStr := strconv.Itoa(index) + "."
	println(indexStr + img.Ou + " -> DONE")
	wg.Done()
}

func GetImageLinks(term string, imageType string, index int) []string {
	url := "https://www.bing.com/images/search?scenario=ImageBasicHover&datsrc=N_I&layout=RowBased&mmasync=1"
	url += "&q=" + term + "&pq=" + term

	imageType = strings.Trim(imageType, " ")
	if len(imageType) == 0 || imageType != "*" {
		url += "&filetype=" + imageType
	}

	url += "&first=" + strconv.Itoa(index) + "&count=35&relp=35"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	page := string(body)

	println(page)

	r := regexp.MustCompile("<div class=\"rg_meta notranslate\">([\\s\\S]*?)</div>")
	imageLinks := r.FindAllStringSubmatch(page, -1)

	var result []string

	for _, ImageItem := range imageLinks {
		result = append(result, ImageItem[1])
	}

	return result
}

func GetFileFullName(img im.ImageItem, folderPath string) string {
	url := img.Ou
	fileName := img.ID[0 : len(img.ID)-1]

	if len(img.Ity) != 0 {
		fileName += "_" + url[strings.LastIndex(img.Ou, "/")+1:strings.LastIndex(img.Ou, ".")] + "." + img.Ity
	} else {
		fileName += ".jpeg"
	}

	return folderPath + "\\" + im.CleanFileName(fileName)
}

func GetImageItemFromJson(jsonString string) (im.ImageItem, error) {
	img := im.ImageItem{}

	err := json.Unmarshal([]byte(jsonString), &img)
	if err != nil {
		return img, err
	}

	return img, nil
}