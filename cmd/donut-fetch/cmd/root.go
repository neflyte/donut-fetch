package cmd

import (
	"fmt"

	"github.com/neflyte/donut-fetch/internal"
	"github.com/spf13/cobra"
)

const (
	defaultFetchTimeout = uint(5)
)

var (
	// AppVersion is the application version number
	AppVersion = "dev"

	// rootCmd is the application's root command
	rootCmd = &cobra.Command{
		Use:   "donut-fetch <sources.json>",
		Short: "Fetch hostname sources for donutdns",
		// TODO: Fetch, if newer, each hostname list...
		Long: "Fetch each hostname list from sources.json and combine into a single list",
		Args: cobra.ExactArgs(1),
		RunE: doRun,
	}
	fetchTimeout uint
	outputFile   string
	numFetches   uint
)

func init() {
	rootCmd.Flags().UintVar(&fetchTimeout, "timeout", defaultFetchTimeout, "connection timeout in seconds")
	rootCmd.Flags().StringVar(&outputFile, "output", "", "output file; output to console if not specified")
	rootCmd.Flags().UintVar(&numFetches, "fetches", 5, "the number of simultaneous downloads")
	rootCmd.SetVersionTemplate(fmt.Sprintf("donut-fetch %s", AppVersion))
}

func Execute() error {
	return rootCmd.Execute()
}

func doRun(_ *cobra.Command, args []string) error {
	log := internal.GetLogger("doRun")
	internal.DefaultFetch.SetTimeout(fetchTimeout)
	sources := internal.NewSources()
	err := sources.FromFile(args[0])
	if err != nil {
		log.Printf("error reading sources file: %s\n", err)
		return err
	}
	state := make(internal.State)
	err = state.Load()
	if err != nil {
		log.Printf("error loading state file: %s\n", err)
		return err
	}
	hosts := internal.NewHosts()
	err = internal.ProcessSites(hosts, sources, &state, numFetches)
	if err != nil {
		log.Printf("error processing sites: %s\n", err)
		return err
	}
	err = state.Save()
	if err != nil {
		log.Printf("error saving state: %s\n", err)
	}
	if outputFile != "" {
		err = hosts.DumpToFile(outputFile)
		if err != nil {
			log.Printf("error dumping hosts to file %s: %s\n", outputFile, err)
			return err
		}
	} else {
		hosts.DumpToConsole()
	}
	return nil
}
