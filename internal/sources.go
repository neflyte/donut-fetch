package internal

import (
	"encoding/json"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
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
	filenameLower := strings.ToLower(filename)
	if strings.HasSuffix(filenameLower, ".yaml") || strings.HasSuffix(filenameLower, ".yml") {
		err = yaml.Unmarshal(fileBytes, &s.sourcesMap)
	} else {
		err = json.Unmarshal(fileBytes, &s.sourcesMap)
	}
	if err != nil {
		log.Printf("error unmarshaling sources: %s\n", err)
		return err
	}
	return nil
}
