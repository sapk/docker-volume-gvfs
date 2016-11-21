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
== Version: %s - Built: %s ==
`
)

var (
	Version   string = ""
	BuildDate string = ""
	rootCmd          = &cobra.Command{
		Use:              "docker-volume-gvfs",
		Short:            "GVfs - Docker volume driver plugin",
		Long:             longHelp,
		PersistentPreRun: setupLogger,
	}
)

func Start() {
	setupFlags()
	rootCmd.Long = fmt.Sprintf(longHelp, Version, BuildDate)
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
