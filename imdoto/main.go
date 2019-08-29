package main

import google "github.com/andrei-galkin/imdoto/google_search"
import bing "github.com/andrei-galkin/imdoto/bing_search"
import yandex "github.com/andrei-galkin/imdoto/yandex_search"
import shared "github.com/andrei-galkin/imdoto/shared"

func main() {
	setting := shared.GetSetting()

	switch setting.Engine {
	case "google":
		google.Download(setting)
	case "bing":
		bing.Download(setting)
	case "yandex":
		yandex.Download(setting)
	}
}
