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
)

import im "github.com/andrei-galkin/imdoto/imdoto"

type ImageItem struct {
	Cnt    string `json:"cnt"`
	Blocks []struct {
		Name struct {
			Block string `json:"block"`
			Mods  struct {
			} `json:"mods"`
		} `json:"name"`
		Params struct {
			PageNum int           `json:"pageNum"`
			Bundles []interface{} `json:"bundles"`
		} `json:"params"`
		HTML string `json:"html"`
	} `json:"blocks"`
	Metadata struct {
		BundlesMetadata struct {
			Lb string `json:"lb"`
		} `json:"bundlesMetadata"`
		AssetsMetadata struct {
			Las string `json:"las"`
		} `json:"assetsMetadata"`
	} `json:"metadata"`
	Assets struct {
		Assets []interface{} `json:"assets"`
	} `json:"assets"`
}

var wg sync.WaitGroup

func Download(option im.Setting) {
	var imageLinks []string
	imageIndex := 0

	for index := 1; index <= option.Limit; index++ {
		if imageIndex == 0 {
			imageLinks = GetImageLinks(option.Term, option.ImageType, index)
		}

		// img, err := GetImageItemFromJson(imageLinks[imageIndex])
		// if err != nil {
		// 	im.PrintError(err)
		// }

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
	fullName := GetFileFullName(url, folderPath)

	if err := DownloadFile(fullName, url); err != nil {
		im.PrintError(err)
	}
	indexStr := strconv.Itoa(index) + "."

	println(indexStr + url)
	println(fullName)
	println("DONE")

	wg.Done()
}

func GetImageLinks(term string, imageType string, index int) []string {
	url := "https://yandex.ru/images/search?rpt=image&format=json&text=" + term

	url += "&p=" + strconv.Itoa(index/30)

	imageType = strings.Trim(imageType, " ")
	if len(imageType) == 0 || imageType != "*" {
		url += "&type=" + imageType
	}

	url += "&request={%22blocks%22:[{%22block%22:%22gallery__items:ajax%22,%22params%22:{},%22version%22:2}]}"

	println(url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		im.PrintError(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	page := string(body)
	println(page)
	r := regexp.MustCompile(`img_url=([\s\S]*?)&amp;text=`)
	imageLinks := r.FindAllStringSubmatch(page, -1)

	var result []string

	for _, ImageItem := range imageLinks {
		result = append(result, strings.Replace(ImageItem[1], "%3A", ":", -1))
	}

	return result
}

func GetFileFullName(img string, folderPath string) string {
	url := img
	fileName := url

	if strings.LastIndex(url, "?") != -1 {
		fileName = url[strings.LastIndex(url, "/")+1 : strings.LastIndex(url, "?")]
	}

	if strings.LastIndex(url, ".") == -1 {
		fileName += ".jpeg"
	}

	return folderPath + "\\" + im.CleanFileName(fileName)
}

func GetImageItemFromJson(jsonString string) (ImageItem, error) {
	img := ImageItem{}

	err := json.Unmarshal([]byte(jsonString), &img)
	if err != nil {
		return img, err
	}

	return img, nil
}
