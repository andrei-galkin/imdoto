package main

import (
	"log"

	search_engine "github.com/andrei-galkin/imdoto/search_engine"
	shared "github.com/andrei-galkin/imdoto/shared"
)

func main() {
	setting, err := shared.GetSetting()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	if err := shared.ValidateSetting(setting); err != nil {
		log.Fatalf("Invalid settings: %v", err)
	}

	s, err := search_engine.NewSearchEngine(setting.Engine)
	if err != nil {
		log.Fatalf("Invalid engine: %v", err)
	}

	if err := s.Download(setting); err != nil {
		log.Fatalf("Download failed: %v", err)
	}
}
