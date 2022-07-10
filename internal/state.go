package internal

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

const (
	stateFileName = "state.json"
	stateFileMode = 0640
	stateDirMode  = 0750
)

type State map[string]ResourceState

func (st *State) Load() error {
	log := GetLogger("Load")
	confDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("error getting user config directory: %s\n", err)
		return err
	}
	stateDir := path.Join(confDir, "donut-fetch")
	err = os.MkdirAll(stateDir, stateDirMode)
	if err != nil {
		log.Printf("error creating state directory: %s\n", err)
		return err
	}
	stateFile := path.Join(stateDir, stateFileName)
	stateBytes, err := os.ReadFile(stateFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		log.Printf("error reading state file: %s\n", err)
		return err
	}
	err = json.Unmarshal(stateBytes, st)
	if err != nil {
		log.Printf("error unmarshalling state from JSON: %s\n", err)
		return err
	}
	return nil
}

func (st *State) Save() error {
	log := GetLogger("Save")
	confDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("error getting user config directory: %s\n", err)
		return err
	}
	stateDir := path.Join(confDir, "donut-fetch")
	err = os.MkdirAll(stateDir, stateDirMode)
	if err != nil {
		log.Printf("error creating state directory: %s\n", err)
		return err
	}
	stateFile := path.Join(stateDir, stateFileName)
	stateBytes, err := json.Marshal(st)
	if err != nil {
		log.Printf("error marshalling state to JSON: %s\n", err)
		return err
	}
	err = os.WriteFile(stateFile, stateBytes, stateFileMode)
	if err != nil {
		log.Printf("error writing state to state file: %s\n", err)
		return err
	}
	return nil
}

func (st *State) Get(url string) ResourceState {
	s, ok := (*st)[HashURL(url)]
	if !ok {
		return ResourceState{}
	}
	return s
}

func (st *State) Set(url string, newState ResourceState) {
	(*st)[HashURL(url)] = newState
}
