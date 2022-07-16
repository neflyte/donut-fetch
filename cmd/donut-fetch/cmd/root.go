package cmd

import (
	"fmt"

	"github.com/neflyte/donut-fetch/internal"
	"github.com/spf13/cobra"
)

const (
	defaultFetchTimeout = uint(5)
	defaultNumFetches   = 5
)

var (
	// AppVersion is the application version number
	AppVersion = "dev"

	// rootCmd is the application's root command
	rootCmd = &cobra.Command{
		Version: AppVersion,
		Use:     "donut-fetch <sources.json|yaml>",
		Short:   "Fetch hostname sources for donutdns",
		Long:    "Fetch, if newer, each hostname list from the sources file and combine into a single list",
		Args:    cobra.ExactArgs(1),
		RunE:    doRun,
	}
	fetchTimeout uint
	outputFile   string
	numFetches   uint
)

func init() {
	rootCmd.Flags().UintVar(&fetchTimeout, "timeout", defaultFetchTimeout, "connection timeout in seconds")
	rootCmd.Flags().StringVar(&outputFile, "output", "", "output file; output to console if not specified")
	rootCmd.Flags().UintVar(&numFetches, "fetches", defaultNumFetches, "the number of simultaneous downloads")
	rootCmd.Flags().BoolVar(&internal.Debug, "debug", false, "enable debug logging")
	rootCmd.SetVersionTemplate(fmt.Sprintf("donut-fetch %s\n", AppVersion))
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
		log.Printf("Wrote %d hosts to file %s\n", hosts.Len(), outputFile)
	} else {
		hosts.DumpToConsole()
	}
	log.Println("done.")
	return nil
}
