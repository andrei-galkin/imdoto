package yandexsearch

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

	shared "github.com/andrei-galkin/imdoto/shared"
)

type ImageItem struct {
	ID   string `json:"id"`
	Dups []struct {
		URL             string `json:"url"`
		FileSizeInBytes int    `json:"fileSizeInBytes"`
		W               int    `json:"w"`
		H               int    `json:"h"`
		Origin          struct {
			W   int    `json:"w"`
			H   int    `json:"h"`
			URL string `json:"url"`
		} `json:"origin"`
		IsMixedImage bool `json:"isMixedImage"`
	} `json:"dups"`
}

var wg sync.WaitGroup

func Download(option shared.Setting) {
	var imageLinks []string
	imageIndex := 0

	for index := 1; index <= option.Limit; index++ {
		if imageIndex == 0 {
			imageLinks = GetImageLinks(option.Term, option.ImageType, index)
		}

		wg.Add(1)
		go DownloadImage(imageLinks[imageIndex], option.FolderPath, index)

		imageIndex += 1

		//exit if there is less images then limit
		if imageIndex == len(imageLinks)-1 && len(imageLinks) != 30 {
			break
		}

		if imageIndex == 30 {
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

func DownloadImage(url string, folderPath string, index int) {
	fullName := GetFileFullNameFromURL(url, folderPath)

	if err := DownloadFile(fullName, url); err != nil {
		shared.PrintError(err)
	}
	indexStr := strconv.Itoa(index) + "."

	println(indexStr + url)
	println(fullName)
	println("DONE")

	wg.Done()
}

// GetFileFullNameFromURL builds a filename when only the image URL is available.
func GetFileFullNameFromURL(url string, folderPath string) string {
	fileName := ""

	if strings.LastIndex(url, "/") != -1 {
		fileName = url[strings.LastIndex(url, "/")+1 : len(url)]
	}

	if strings.LastIndex(url, "?") != -1 {
		fileName = url[strings.LastIndex(url, "/")+1 : strings.LastIndex(url, "?")]
	}

	if strings.LastIndex(fileName, ".") == -1 {
		fileName += ".jpeg"
	}

	return folderPath + "\\" + shared.CleanFileName(fileName)
}

func GetImageLinks(term string, imageType string, index int) []string {
	url := "https://yandex.ru/images/search?rpt=image&format=json&text=" + term

	url += "&p=" + strconv.Itoa(index/30)

	imageType = strings.Trim(imageType, " ")
	if len(imageType) == 0 || imageType != "*" {
		url += "&type=" + imageType
	}

	url += "&request={%22blocks%22:[{%22block%22:%22gallery__items:ajax%22,%22params%22:{},%22version%22:2}]}"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		shared.PrintError(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	page := string(body)

	r := regexp.MustCompile(`img_url=([\s\S]*?)&amp;text=`)
	imageLinks := r.FindAllStringSubmatch(page, -1)

	var result []string

	for _, ImageItem := range imageLinks {
		result = append(result, strings.Replace(ImageItem[1], "%3A", ":", -1))
	}

	return result
}

func GetFileFullName(img ImageItem, folderPath string) string {
	url := ""
	fileName := ""

	if len(img.Dups) > 0 {
		url = img.Dups[0].Origin.URL
	}

	if url == "" {
		return folderPath + "\\" + shared.CleanFileName(fileName)
	}

	if strings.LastIndex(url, "/") != -1 {
		fileName = url[strings.LastIndex(url, "/")+1 : len(url)]
	}

	if strings.LastIndex(url, "?") != -1 {
		fileName = url[strings.LastIndex(url, "/")+1 : strings.LastIndex(url, "?")]
	}

	if len(img.ID) != 0 {
		fileName = img.ID + "_" + fileName
	}

	if strings.LastIndex(fileName, ".") == -1 {
		fileName += ".jpeg"
	}

	return folderPath + "\\" + shared.CleanFileName(fileName)
}

func GetImageItemFromJson(jsonString string) (ImageItem, error) {
	img := ImageItem{}

	err := json.Unmarshal([]byte(jsonString), &img)
	if err != nil {
		return img, err
	}

	return img, nil
}
