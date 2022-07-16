package internal

import (
	"errors"
	"strings"
	"sync"

	"github.com/panjf2000/ants/v2"
)

func ProcessSites(hosts *Hosts, sources *Sources, state *State, numFetches uint) error {
	wg := sync.WaitGroup{}
	log := GetLogger("ProcessSites")
	pool, err := ants.NewPool(int(numFetches))
	if err != nil {
		log.Printf("error creating new pool: %s\n", err)
		return err
	}
	defer pool.Release()
	for category, sites := range sources.Sources() {
		for x := range sites {
			site := sites[x]
			err := pool.Submit(func() {
				processSiteWithOptions(processSiteOptions{
					hosts: hosts,
					state: state,
					site:  &site,
					wg:    &wg,
				})
			})
			if err != nil {
				log.Printf("error processing site %s of category %s: %s\n", sites[x], category, err)
				continue
			}
			wg.Add(1)
		}
	}
	log.Println("waiting for processing to complete")
	wg.Wait()
	log.Println("processing completed")
	return nil
}

type processSiteOptions struct {
	hosts *Hosts
	state *State
	site  *string
	wg    *sync.WaitGroup
}

func processSiteWithOptions(options processSiteOptions) {
	log := GetLogger("processSiteWithOptions")
	defer options.wg.Done()
	err := processSite(options.hosts, options.state, options.site)
	if err != nil {
		log.Printf("error processing site %s: %s\n", *options.site, err)
	}
}

func processSite(hosts *Hosts, state *State, siteURL *string) error {
	var (
		err     error
		fetched []string
	)

	log := GetLogger("processSite")
	if siteURL == nil {
		return errors.New("unexpected nil site url")
	}
	site := *siteURL
	siteLocalState := state.Get(site)
	// get the state of the resource
	siteState, err := DefaultFetch.ResourceState(site)
	if err != nil {
		log.Printf("error getting resource state of site: %s\n", err)
		return err
	}
	cacheID := HashURL(strings.ToLower(site))
	if siteLocalState.IsETagStale(siteState.ETag) || siteLocalState.IsLastModifiedPast(siteState.LastModified) {
		log.Printf("fetching new hosts from %s\n", site)
		fetched, err = DefaultFetch.Hosts(site)
		if err != nil {
			log.Printf("error fetching site %s: %s\n", site, err)
			return err
		}
		// update cached hosts with fetched
		cacheHosts := NewHosts()
		cacheHosts.Add(fetched)
		err = cacheHosts.SaveCache(cacheID)
		if err != nil {
			log.Printf("error saving hosts to cache: %s\n", err)
		}
		hosts.Add(fetched)
		state.Set(site, siteState)
	} else {
		// load cached hosts and add to list
		log.Printf("loading hosts from cache for %s\n", site)
		cachedHosts := NewHosts()
		err = cachedHosts.LoadCache(cacheID)
		if err != nil {
			log.Printf("error loading hosts from cache: %s\n", err)
			return err
		}
		hosts.Add(cachedHosts.Hosts())
	}
	return nil
}
