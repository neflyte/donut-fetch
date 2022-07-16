package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"golang.org/x/exp/maps"
)

const (
	dumpfileMode  = 0640
	cacheFileMode = 0640
	cacheDirMode  = 0750
)

type Hosts struct {
	hostnames map[string]bool
	maplock   sync.RWMutex
}

func NewHosts() *Hosts {
	return &Hosts{
		hostnames: make(map[string]bool),
		maplock:   sync.RWMutex{},
	}
}

func (h *Hosts) Clear() {
	h.maplock.Lock()
	defer h.maplock.Unlock()
	maps.Clear(h.hostnames)
}

func (h *Hosts) Len() uint {
	h.maplock.RLock()
	defer h.maplock.RUnlock()
	return uint(len(h.hostnames))
}

func (h *Hosts) Hosts() []string {
	h.maplock.RLock()
	defer h.maplock.RUnlock()
	return maps.Keys(h.hostnames)
}

func (h *Hosts) Add(newHosts []string) {
	h.maplock.Lock()
	defer h.maplock.Unlock()
	for x := range newHosts {
		h.hostnames[newHosts[x]] = true
	}
}

func (h *Hosts) DumpToConsole() {
	fmt.Printf("%d hosts\n", h.Len())
	for _, hostname := range h.Hosts() {
		fmt.Println(hostname)
	}
}

func (h *Hosts) DumpToFile(filename string) error {
	log := GetLogger("DumpToFile")
	fout, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, dumpfileMode)
	if err != nil {
		log.Printf("error opening dump file: %s\n", err)
		return err
	}
	defer func() {
		err = fout.Close()
		if err != nil {
			log.Printf("error closing dump file: %s\n", err)
		}
	}()
	for _, host := range h.Hosts() {
		_, err = fout.WriteString(host + "\n")
		if err != nil {
			log.Printf("error writing to dump file: %s\n", err)
			return err
		}
	}
	return nil
}

func (h *Hosts) LoadCache(cacheID string) error {
	log := GetLogger("LoadCache")
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Printf("error getting user cache directory: %s\n", err)
		return err
	}
	hostsCacheDir := path.Join(userCacheDir, "donut-fetch")
	err = os.MkdirAll(hostsCacheDir, cacheDirMode)
	if err != nil {
		log.Printf("error creating cache directory %s: %s\n", hostsCacheDir, err)
		return err
	}
	cacheFile := path.Join(hostsCacheDir, cacheID+".json")
	cacheBytes, err := os.ReadFile(cacheFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		log.Printf("error reading cache file %s: %s\n", cacheFile, err)
		return err
	}
	cacheHosts := make([]string, 0)
	err = json.Unmarshal(cacheBytes, &cacheHosts)
	if err != nil {
		log.Printf("error unmarshalling cache data from JSON: %s\n", err)
		return err
	}
	h.Add(cacheHosts)
	log.Printf("loaded %d hosts from cache file %s\n", len(cacheHosts), cacheFile)
	return nil
}

func (h *Hosts) SaveCache(cacheID string) error {
	log := GetLogger("SaveCache")
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Printf("error getting user cache directory: %s\n", err)
		return err
	}
	hostsCacheDir := path.Join(userCacheDir, "donut-fetch")
	err = os.MkdirAll(hostsCacheDir, cacheDirMode)
	if err != nil {
		log.Printf("error creating cache directory %s: %s\n", hostsCacheDir, err)
		return err
	}
	cacheFile := path.Join(hostsCacheDir, cacheID+".json")
	cacheBytes, err := json.Marshal(h.Hosts())
	if err != nil {
		log.Printf("error marshalling hostnames to JSON: %s\n", err)
		return err
	}
	err = os.WriteFile(cacheFile, cacheBytes, cacheFileMode)
	if err != nil {
		log.Printf("error writing cache file %s: %s\n", cacheFile, err)
		return err
	}
	log.Printf("wrote %d hosts to cache file %s\n", h.Len(), cacheFile)
	return nil
}
