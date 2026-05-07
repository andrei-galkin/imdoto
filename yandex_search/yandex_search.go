package yandexsearch

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
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

func NewSearchEngine() *YandexSearch {
	return &YandexSearch{}
}

type YandexSearch struct{}

func (ys *YandexSearch) Download(option shared.Setting) error {
	var imageLinks []string
	imageIndex := 0
	sem := make(chan struct{}, shared.MaxConcurrentDownloads)

	for index := 1; index <= option.Limit; index++ {
		if imageIndex == 0 {
			var err error
			imageLinks, err = GetImageLinks(option.Term, option.ImageType, index)
			if err != nil {
				return shared.WrapError(err, "failed to fetch image links")
			}
			if len(imageLinks) == 0 {
				return errors.New("no image links found")
			}
		}

		wg.Add(1)
		go DownloadImage(imageLinks[imageIndex], option.FolderPath, index, sem)

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
	return nil
}

func DownloadImage(url string, folderPath string, index int, sem chan struct{}) {
	sem <- struct{}{}
	defer func() { <-sem }()

	fullName := GetFileFullNameFromURL(url, folderPath)

	if err := shared.DownloadFile(fullName, url); err != nil {
		shared.LogError(shared.WrapError(err, "failed to download image"))
	}
	indexStr := strconv.Itoa(index) + "."

	println(indexStr + url)
	println(fullName)
	println("DONE")

	wg.Done()
}

// GetFileFullNameFromURL builds a filename when only the image URL is available.
func GetFileFullNameFromURL(url string, folderPath string) string {
	folderPath = shared.NormalizeFolderPath(folderPath)

	urlBase := strings.SplitN(url, "?", 2)[0]
	fileName := filepath.Base(urlBase)

	if strings.LastIndex(fileName, ".") == -1 {
		fileName += ".jpeg"
	}

	return filepath.Join(folderPath, shared.CleanFileName(fileName))
}

func GetImageLinks(term string, imageType string, index int) ([]string, error) {
	url := "https://yandex.ru/images/search?rpt=image&format=json&text=" + term

	url += "&p=" + strconv.Itoa(index/30)

	imageType = strings.Trim(imageType, " ")
	if len(imageType) == 0 || imageType != "*" {
		url += "&type=" + imageType
	}

	url += "&request={%22blocks%22:[{%22block%22:%22gallery__items:ajax%22,%22params%22:{},%22version%22:2}]}"

	client := shared.HTTPClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, shared.WrapError(err, "failed to create HTTP request")
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, shared.WrapError(err, "failed to fetch image results")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, shared.WrapError(err, "failed to read response body")
	}
	page := string(body)

	r := regexp.MustCompile(`img_url=([\s\S]*?)&amp;text=`)
	imageLinks := r.FindAllStringSubmatch(page, -1)

	var result []string
	for _, ImageItem := range imageLinks {
		result = append(result, strings.Replace(ImageItem[1], "%3A", ":", -1))
	}

	return result, nil
}

func GetFileFullName(img ImageItem, folderPath string) string {
	folderPath = shared.NormalizeFolderPath(folderPath)
	url := ""
	fileName := ""

	if len(img.Dups) > 0 {
		url = img.Dups[0].Origin.URL
	}

	if url == "" {
		return filepath.Join(folderPath, shared.CleanFileName(fileName))
	}

	if strings.LastIndex(url, "/") != -1 {
		fileName = url[strings.LastIndex(url, "/")+1:]
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

	return filepath.Join(folderPath, shared.CleanFileName(fileName))
}

func GetImageItemFromJson(jsonString string) (ImageItem, error) {
	img := ImageItem{}

	err := json.Unmarshal([]byte(jsonString), &img)
	if err != nil {
		return img, err
	}

	return img, nil
}
