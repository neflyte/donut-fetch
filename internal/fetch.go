package internal

import (
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	DefaultFetch = NewFetch()
)

type Fetch struct {
	client http.Client
}

func NewFetch() *Fetch {
	return &Fetch{
		client: http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: -1,
				DisableKeepAlives:   true,
			},
		},
	}
}

func (f *Fetch) SetTimeout(timeout uint) {
	f.client.Timeout = time.Duration(timeout) * time.Second
}

func (f *Fetch) Hosts(url string) ([]string, error) {
	log := GetLogger("Hosts")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("error creating get request: %s\n", err)
		return nil, err
	}
	res, err := f.client.Do(req)
	if err != nil {
		log.Printf("error executing get request: %s\n", err)
		return nil, err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			log.Printf("error closing response body: %s\n", err)
		}
	}()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error reading response body: %s\n", err)
		return nil, err
	}
	rawtoks := strings.Split(string(bodyBytes), "\n")
	hosts := make([]string, 0)
	for x := range rawtoks {
		trimmed := strings.TrimSpace(rawtoks[x])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		// split on space in case it's in hosts format
		hostToks := strings.Split(trimmed, " ")
		if len(hostToks) > 1 {
			trimmed = hostToks[1]
		}
		// split on hash in case there are mid-line comments
		commentToks := strings.Split(trimmed, "#")
		if len(commentToks) > 1 {
			trimmed = commentToks[0]
		}
		// add to list of hosts
		hosts = append(hosts, trimmed)
	}
	return hosts, nil
}

func (f *Fetch) ResourceState(url string) (ResourceState, error) {
	log := GetLogger("LastModified")
	rState := ResourceState{}
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		log.Printf("error creating get request: %s\n", err)
		return rState, err
	}
	res, err := f.client.Do(req)
	if err != nil {
		log.Printf("error executing get request: %s\n", err)
		return rState, err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			log.Printf("error closing response body: %s\n", err)
		}
	}()
	lastModHeader := res.Header.Get("Last-Modified")
	if lastModHeader != "" {
		lastModTime, err := time.Parse(time.RFC1123, lastModHeader)
		if err != nil {
			log.Printf("error parsing last-modified header: %s\n", err)
		} else {
			rState.LastModified = lastModTime
			log.Printf("lastMod=%s for url %s\n", lastModTime.String(), url)
		}
	}
	rState.ETag = res.Header.Get("ETag")
	return rState, nil
}
