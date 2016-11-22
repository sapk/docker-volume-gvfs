package driver

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	VerboseFlag = "verbose"
	longHelp    = `
docker-volume-gvfs (GVfs Volume Driver Plugin)
Provides docker volume support for GVfs.
== Version: %s - Commit: %s ==
`
)

var (
	Version string
	Commit  string
	rootCmd = &cobra.Command{
		Use:              "docker-volume-gvfs",
		Short:            "GVfs - Docker volume driver plugin",
		Long:             longHelp,
		PersistentPreRun: setupLogger,
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display current version and build date",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nVersion: %s - Commit: %s\n\n", Version, Commit)
		},
	}
)

func Start() {
	setupFlags()
	rootCmd.Long = fmt.Sprintf(longHelp, Version, Commit)
	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()
}

func setupFlags() {

}

func setupLogger(cmd *cobra.Command, args []string) {
	if verbose, _ := cmd.Flags().GetBool(VerboseFlag); verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
