package driver

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/spf13/cobra"
)

const (
	VerboseFlag = "verbose"
	BasedirFlag = "basedir"
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
	daemonCmd = &cobra.Command{
		Use:   "daemon",
		Short: "Run plugin in deamon mode",
		Run:   daemonStart,
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display current version and build date",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nVersion: %s - Commit: %s\n\n", Version, Commit)
		},
	}
	baseDir     = ""
	PluginAlias = "gvfs"
)

func Start() {
	setupFlags()
	rootCmd.Long = fmt.Sprintf(longHelp, Version, Commit)
	rootCmd.AddCommand(versionCmd, daemonCmd)
	rootCmd.Execute()
}
func daemonStart(cmd *cobra.Command, args []string) {
	//TODO get args
	//TODO support -o of gvfsd-fuse
	driver := newGVfsDriver(baseDir)
	//newGVfsDriver(baseDir)
	log.Debug(driver)
	h := volume.NewHandler(driver) //TODO subscribe and handle
	log.Debug(h)
	err := h.ServeUnix("root", "gvfs") //TODO test
	if err != nil {
		log.Debug(err)
	}
}

func setupFlags() {
	rootCmd.PersistentFlags().StringVar(&baseDir, BasedirFlag, filepath.Join(volume.DefaultDockerRootDirectory, PluginAlias), "Mounted volume base directory")
}

func setupLogger(cmd *cobra.Command, args []string) {
	if verbose, _ := cmd.Flags().GetBool(VerboseFlag); verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
