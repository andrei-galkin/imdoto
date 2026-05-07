package googlesearch

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
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
		go DownloadImage(img, setting.FolderPath, index)

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

func DownloadImage(img ImageItem, folderPath string, index int) {
	fullName := GetFileFullName(img, folderPath)

	if err := DownloadFile(fullName, img.Ou); err != nil {
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

	client := &http.Client{}
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
	url := img.Ou
	fileName := ""

	if len(img.ID) > 1 {
		fileName = img.ID[0 : len(img.ID)-1]
	} else if len(img.ID) == 1 {
		fileName = img.ID
	}

	if len(img.Ity) != 0 {
		start := strings.LastIndex(img.Ou, "/")
		end := strings.LastIndex(img.Ou, ".")
		if start != -1 && end != -1 && end > start {
			fileName += "_" + url[start+1:end] + "." + img.Ity
		} else {
			fileName += "." + img.Ity
		}
	} else {
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
