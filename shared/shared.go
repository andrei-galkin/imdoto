package shared

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Setting struct {
	Engine     string
	Term       string
	FolderName string
	Limit      int
	ImageType  string
	FolderPath string
}

func GetSetting() (Setting, error) {
	engine := flag.String("engine", "yandex", "search engine: google, bing, or yandex")
	folderName := flag.String("folder", "img", "destination folder name")
	term := flag.String("term", "apple", "search term")
	limit := flag.Int("limit", 75, "maximum number of images to download")
	imageType := flag.String("type", "*", "image file type filter (e.g., jpeg, png, or *)")
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return Setting{}, err
	}

	folderPath := filepath.Join(dir, *folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
			return Setting{}, err
		}
	}

	var setting Setting
	setting.Term = strings.Replace(*term, " ", "+", -1)
	setting.FolderName = *folderName
	setting.FolderPath = folderPath
	setting.Limit = *limit
	setting.ImageType = *imageType
	setting.Engine = *engine

	return setting, nil
}

func ValidateSetting(setting Setting) error {
	switch setting.Engine {
	case "google", "bing", "yandex":
	default:
		return fmt.Errorf("unsupported engine: %q", setting.Engine)
	}

	if len(setting.Term) == 0 {
		return fmt.Errorf("search term cannot be empty")
	}

	if setting.Limit <= 0 {
		return fmt.Errorf("limit must be greater than zero")
	}

	if len(setting.FolderName) == 0 {
		return fmt.Errorf("folder name cannot be empty")
	}

	if len(setting.ImageType) == 0 {
		return fmt.Errorf("image type cannot be empty")
	}

	return nil
}

var HTTPClient = &http.Client{
	Timeout: 30 * time.Second,
}

func DownloadFile(filePath string, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func NormalizeFolderPath(folderPath string) string {
	if len(folderPath) == 2 && folderPath[1] == ':' {
		return folderPath + string(os.PathSeparator)
	}
	return filepath.Clean(folderPath)
}

func CleanFileName(fileName string) string {
	symbols := [10]string{"*", "?", "%", "\\", "/", " ", "+", "#", "@", "~"}
	for _, symbol := range symbols {
		fileName = strings.Replace(fileName, symbol, "", -1)
	}
	return fileName
}

func LogError(err error) {
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
}

func WrapError(err error, context string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", context, err)
	}
	return nil
}

func PrintError(err error) {
	LogError(err)
}
