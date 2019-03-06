package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
)

type logWriter struct{}

func main() {
	// fileUrl := "https://golangcode.com/images/avatar.jpg"

	// if err := DownloadImage("image.jpg", fileUrl); err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Success!")

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
	fmt.Println(reflect.TypeOf(imageSource))

	r := regexp.MustCompile("<div class=\"rg_meta notranslate\">([\\s\\S]*?)</div>")

	matches := r.FindAllStringSubmatch(imageSource, -1)
	fmt.Println("GOT : ", len(matches))

	for index, each := range matches {
		fmt.Println("--------->\n", index, each)
	}
}

func DownloadImage(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
