package internal

import (
	"encoding/json"
	"os"
)

type Sources struct {
	sourcesMap map[string][]string
}

func NewSources() *Sources {
	return &Sources{
		sourcesMap: make(map[string][]string),
	}
}

func (s *Sources) Sources() map[string][]string {
	return s.sourcesMap
}

func (s *Sources) FromFile(filename string) error {
	log := GetLogger("FromFile")
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("error reading file %s: %s\n", filename, err)
		return err
	}
	err = json.Unmarshal(fileBytes, &s.sourcesMap)
	if err != nil {
		log.Printf("error unmarshaling sources from JSON: %s\n", err)
		return err
	}
	return nil
}
