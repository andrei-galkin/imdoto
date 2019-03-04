package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	fileUrl := "https://golangcode.com/images/avatar.jpg"

	if err := DownloadImage("image.jpg", fileUrl); err != nil {
		panic(err)
	}

	fmt.Println("Success!")
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
