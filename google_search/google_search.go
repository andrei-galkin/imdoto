package googlesearch

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

type GoogleSearch struct{}

var wg sync.WaitGroup

func NewSearchEngine() *GoogleSearch {
	return &GoogleSearch{}
}

func (gs *GoogleSearch) Download(setting shared.Setting) error {
	var imageLinks []string
	imageIndex := 0
	sem := make(chan struct{}, shared.MaxConcurrentDownloads)

	for index := 1; index <= setting.Limit; index++ {
		if imageIndex == 0 {
			var err error
			imageLinks, err = GetImageLinks(setting.Term, setting.ImageType, index-1)
			if err != nil {
				return shared.WrapError(err, "failed to fetch image links")
			}
			if len(imageLinks) == 0 {
				return errors.New("no image links found")
			}
		}
		if imageIndex >= len(imageLinks) {
			// no more items in current batch; try resetting to fetch next batch
			imageIndex = 0
			var err error
			imageLinks, err = GetImageLinks(setting.Term, setting.ImageType, index)
			if err != nil {
				return shared.WrapError(err, "failed to fetch next batch of image links")
			}
			if len(imageLinks) == 0 {
				return errors.New("no image links found")
			}
		}

		img, err := GetImageItemFromJson(imageLinks[imageIndex])
		if err != nil {
			return err
		}

		imageIndex += 1

		wg.Add(1)
		go DownloadImage(img, setting.FolderPath, index, sem)

		//exit if there is less images then limit
		if imageIndex == len(imageLinks)-1 && len(imageLinks) != 100 {
			break
		}

		if imageIndex == 100 {
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

	if err := shared.DownloadFile(fullName, img.Ou); err != nil {
		shared.LogError(shared.WrapError(err, "failed to download image"))
	}
	indexStr := strconv.Itoa(index) + "."

	println(indexStr + img.Ou)
	println(fullName)
	println("DONE")

	wg.Done()
}

func GetImageLinks(term string, imageType string, index int) ([]string, error) {
	url := "https://www.google.com/search?q=" + term + "&oq=" + term

	imageType = strings.Trim(imageType, " ")
	if len(imageType) == 0 || imageType != "*" {
		url += "&tbs=ift:" + imageType
	}

	if index == 0 {
		url += "&biw=1536&bih=723&tbm=isch&sa=1&ei=6qqGXM_oDenYjwSw1b-oAw"
	} else {
		url += "&ijn=" + strconv.Itoa(index/100) + "&start=" + strconv.Itoa(index) +
			"&biw=1536&bih=723&tbm=isch&sa=1&ei=6qqGXM_oDenYjwSw1b-oAw&yv=3&as_st=y&tbm=isch&asearch=ichunk&async=_id:rg_s,_pms:s,_fmt:pc"
	}

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
		return nil, shared.WrapError(err, "failed to read image results response")
	}
	page := string(body)

	r := regexp.MustCompile("<div class=\"rg_meta notranslate\">([\\s\\S]*?)</div>")
	imageLinks := r.FindAllStringSubmatch(page, -1)

	var result []string
	for _, ImageItem := range imageLinks {
		result = append(result, ImageItem[1])
	}

	return result, nil
}

func GetFileFullName(img ImageItem, folderPath string) string {
	folderPath = shared.NormalizeFolderPath(folderPath)

	var fileName string
	if len(img.ID) > 1 {
		fileName = img.ID[:len(img.ID)-1]
	} else if len(img.ID) == 1 {
		fileName = img.ID
	}

	urlBase := strings.SplitN(img.Ou, "?", 2)[0]
	if len(img.Ity) != 0 {
		base := filepath.Base(urlBase)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		if name != "" {
			fileName += "_" + name + "." + img.Ity
		} else {
			fileName += "." + img.Ity
		}
	} else {
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
