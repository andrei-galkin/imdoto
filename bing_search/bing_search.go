package bingsearch

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
	Cid   string `json:"cid"`
	Purl  string `json:"purl"`
	Murl  string `json:"murl"`
	Turl  string `json:"turl"`
	Md5   string `json:"md5"`
	Shkey string `json:"shkey"`
	T     string `json:"t"`
	Mid   string `json:"mid"`
	Desc  string `json:"desc"`
}

var wg sync.WaitGroup

func NewSearchEngine() *BingSearch {
	return &BingSearch{}
}

type BingSearch struct{}

func (bs *BingSearch) Download(option shared.Setting) error {
	var imageLinks []string
	imageIndex := 0
	sem := make(chan struct{}, shared.MaxConcurrentDownloads)

	for index := 1; index <= option.Limit; index++ {
		if imageIndex == 0 {
			var err error
			imageLinks, err = GetImageLinks(option.Term, option.ImageType, index-1)
			if err != nil {
				return shared.WrapError(err, "failed to fetch image links")
			}
			if len(imageLinks) == 0 {
				return errors.New("no image links found")
			}
		}

		img, err := GetImageItemFromJson(imageLinks[imageIndex])
		if err != nil {
			return shared.WrapError(err, "failed to parse image item")
		}

		imageIndex += 1

		wg.Add(1)
		go DownloadImage(img, option.FolderPath, index, sem)

		//exit if there less images then limit
		if imageIndex == len(imageLinks)-1 && len(imageLinks) != 35 {
			break
		}

		if imageIndex == 35 {
			imageIndex = 0
		}
	}

	wg.Wait()
	return nil
}

func DownloadImage(img ImageItem, folderPath string, index int, sem chan struct{}) {
	sem <- struct{}{}
	defer func() { <-sem }()

	fullName := GetFileFullName(img, folderPath)

	if err := shared.DownloadFile(fullName, img.Murl); err != nil {
		shared.LogError(shared.WrapError(err, "failed to download image"))
	}
	indexStr := strconv.Itoa(index) + "."

	println(indexStr + img.Murl)
	println(fullName)
	println("DONE")

	wg.Done()
}

func GetImageLinks(term string, imageType string, index int) ([]string, error) {
	url := "https://www.bing.com/images/async?scenario=ImageBasicHover&datsrc=N_I&layout=RowBased&mmasync=1"
	url += "&q=" + term

	imageType = strings.Trim(imageType, " ")
	if len(imageType) == 0 || imageType != "*" {
		url += "&filetype=" + imageType
	}

	url += "&first=" + strconv.Itoa(index) + "&count=35&relp=35"

	client := shared.HTTPClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, shared.WrapError(err, "failed to create HTTP request")
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")

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

	r := regexp.MustCompile(" m=\"([\\s\\S]*?)\" onclick=\"")
	imageLinks := r.FindAllStringSubmatch(page, -1)

	var result []string
	for _, ImageItem := range imageLinks {
		result = append(result, RestoreQuotes(ImageItem[1]))
	}

	return result, nil
}

func GetFileFullName(img ImageItem, folderPath string) string {
	folderPath = shared.NormalizeFolderPath(folderPath)
	urlBase := strings.SplitN(img.Murl, "?", 2)[0]
	fileName := img.Cid + "_" + filepath.Base(urlBase)

	if strings.LastIndex(fileName, ".") == -1 {
		fileName += ".jpeg"
	}

	if len(fileName) > 250 {
		fileName = img.Cid + fileName[strings.LastIndex(fileName, ".")+1:len(fileName)]
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

func RestoreQuotes(jsonString string) string {
	return strings.Replace(jsonString, "&quot;", "\"", -1)
}
