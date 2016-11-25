package driver

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/spf13/cobra"
)

const (
	//VerboseFlag flag to set more verbose level
	VerboseFlag = "verbose"
	//DBusFlag flag to set DBus path
	DBusFlag = "dbus"
	//EnvDBus env to setor get from session DBus path
	EnvDBus = "DBUS_SESSION_BUS_ADDRESS"
	//BasedirFlag flag to set the basedir of mounted volumes
	BasedirFlag = "basedir"
	longHelp    = `
docker-volume-gvfs (GVfs Volume Driver Plugin)
Provides docker volume support for GVfs.
== Version: %s - Branch: %s - Commit: %s ==
`
)

var (
	//Version version of running code
	Version string
	//Branch branch of running code
	Branch string
	//Commit commit of running code
	Commit string
	//PluginAlias plugin alias name in docker
	PluginAlias = "gvfs"
	baseDir     = ""
	rootCmd     = &cobra.Command{
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
			fmt.Printf("\nVersion: %s - Branch: %s - Commit: %s\n\n", Version, Branch, Commit)
		},
	}
)

//Start start the program
func Start() {
	setupFlags()
	rootCmd.Long = fmt.Sprintf(longHelp, Version, Branch, Commit)
	rootCmd.AddCommand(versionCmd, daemonCmd)
	rootCmd.Execute()
}

func typeOrEnv(cmd *cobra.Command, flag, envname string) string {
	val, _ := cmd.Flags().GetString(flag)
	if val == "" {
		val = os.Getenv(envname)
	}
	return val
}

func daemonStart(cmd *cobra.Command, args []string) {
	//TODO get additional args
	dbus := typeOrEnv(cmd, DBusFlag, EnvDBus)
	driver := newGVfsDriver(baseDir, dbus)
	log.Debug(driver)
	h := volume.NewHandler(driver)
	log.Debug(h)
	err := h.ServeUnix("root", PluginAlias)
	if err != nil {
		log.Debug(err)
	}
}

func setupFlags() {
	rootCmd.PersistentFlags().Bool(VerboseFlag, false, "Turns on verbose logging")
	rootCmd.PersistentFlags().StringVar(&baseDir, BasedirFlag, filepath.Join(volume.DefaultDockerRootDirectory, PluginAlias), "Mounted volume base directory")

	daemonCmd.Flags().StringP(DBusFlag, "d", "", "DBus address to use for gvfs.  Can also set default environment DBUS_SESSION_BUS_ADDRESS")
}

func setupLogger(cmd *cobra.Command, args []string) {
	if verbose, _ := cmd.Flags().GetBool(VerboseFlag); verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
