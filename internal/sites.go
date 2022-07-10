package internal

func ProcessSites(hosts *Hosts, sources *Sources, state *State) error {
	log := GetLogger("ProcessSites")
	for category, sites := range sources.Sources() {
		for x := range sites {
			err := processSite(hosts, state, category, sites[x])
			if err != nil {
				log.Printf("error processing site %s of category %s: %s\n", sites[x], category, err)
				continue
			}
		}
	}
	return nil
}

func processSite(hosts *Hosts, state *State, category, site string) error {
	var (
		err     error
		fetched []string
	)

	log := GetLogger("processSite")
	siteLocalState := state.Get(site)
	// get the state of the resource
	siteState, err := DefaultFetch.ResourceState(site)
	if err != nil {
		log.Printf("error getting resource state of site: %s\n", err)
		return err
	}
	// log.Printf(
	//	"localETag=%s, siteETag=%s; localLastModified=%s, siteLastModified=%s\n",
	//	siteLocalState.ETag,
	//	siteState.ETag,
	//	siteLocalState.LastModified.String(),
	//	siteState.LastModified.String(),
	// )
	cacheID := HashURL(site)
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
		siteLocalState.ETag = siteState.ETag
		siteLocalState.LastModified = siteState.LastModified
		state.Set(site, siteLocalState)
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
