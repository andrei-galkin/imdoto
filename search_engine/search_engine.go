package search_engine

import (
	"fmt"

	bing "github.com/andrei-galkin/imdoto/bing_search"
	google "github.com/andrei-galkin/imdoto/google_search"
	shared "github.com/andrei-galkin/imdoto/shared"
	yandex "github.com/andrei-galkin/imdoto/yandex_search"
)

type SearchEngine interface {
	Download(setting shared.Setting) error
}

func NewSearchEngine(engineName string) (SearchEngine, error) {
	switch engineName {
	case "google":
		return google.NewSearchEngine(), nil
	case "bing":
		return bing.NewSearchEngine(), nil
	case "yandex":
		return yandex.NewSearchEngine(), nil
	default:
		return nil, fmt.Errorf("unknown search engine: %s", engineName)
	}
}
