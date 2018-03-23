package drivers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sapk/docker-volume-helpers/basic"
	"github.com/sapk/docker-volume-helpers/driver"
	"github.com/sirupsen/logrus"
)

const (
	//MountTimeout timeout before killing a mount try in seconds
	MountTimeout = 30
	//CfgVersion current config version compat
	CfgVersion = 2
	//CfgFolder config folder
	CfgFolder = "/etc/docker-volumes/gvfs/"
)

//GVfsDriver docker volume plugin driver extension of basic plugin
type GVfsDriver = basic.Driver

//gvfsVolumeDriver Handle the specific drive module (ssh/dav/ftp,...)
type gvfsVolumeDriver interface {
	id() DriverType
	isAvailable() bool
	mountpoint() (string, error)
}

//Init start all needed deps and serve response to API call
func Init(root string, dbus, fuseOpts string) *GVfsDriver {
	logrus.Debugf("Init gluster driver at %s", root)
	config := basic.DriverConfig{
		Version: CfgVersion,
		Root:    root,
		Folder:  CfgFolder,
		CustomOptions: map[string]interface{}{
			"dbus":     dbus,
			"fuseOpts": fuseOpts,
			"env":      []string{},
		},
	}
	eventHandler := basic.DriverEventHandler{
		OnMountVolume: mountVolume,
		OnInit:        initDriver,
		GetMountName:  GetMountName,
	}
	return basic.Init(&config, &eventHandler)
}

//GetMountName get moint point base on request and driver config
func GetMountName(d *basic.Driver, r *volume.CreateRequest) (string, error) {
	if r.Options == nil || r.Options["url"] == "" {
		return "", fmt.Errorf("url option required")
	}

	_, m, err := getDriver(r.Options["url"])
	if err != nil {
		return "", err
	}
	return filepath.Join(d.Config.Root, m), nil
}

func mountVolume(d *basic.Driver, v driver.Volume, m driver.Mount, r *volume.MountRequest) (*volume.MountResponse, error) {
	opts := v.GetOptions()
	cmd := fmt.Sprintf("gio mount %s", opts["url"])
	if opts["password"] != "" {
		p := setEnv(cmd, d.Config.CustomOptions["env"].([]string))
		inStd, err := p.StdinPipe()
		if err != nil { //Get a input buffer
			return nil, err
		}
		var outStd bytes.Buffer
		p.Stdout = &outStd
		var errStd bytes.Buffer
		p.Stderr = &errStd

		if err := p.Start(); err != nil {
			return nil, err
		}
		inStd.Write([]byte(opts["password"] + "\n")) //Send password to process + Send return line

		// wait or timeout
		donec := make(chan error, 1)
		go func() {
			donec <- p.Wait() //Process finish
		}()
		select {
		case <-time.After(MountTimeout * time.Second):
			sOut := outStd.String()
			sErr := errStd.String()
			p.Process.Kill()
			logrus.Debugf("out : %s", sOut)
			logrus.Debugf("outErr : %s", sErr)
			return nil, fmt.Errorf("The command %s timeout", cmd)
		case <-donec:
			sOut := outStd.String()
			sErr := errStd.String()
			logrus.Debugf("Password send and command %s return", cmd)
			logrus.Debugf("out : %s", sOut)
			logrus.Debugf("outErr : %s", sErr)
			// handle erros like : "Error mounting location: Location is already mounted" or Error mounting location: Could not connect to 10.8.0.7: No route to host
			if strings.Contains(sErr, "Error mounting location") {
				return nil, fmt.Errorf("Error mounting location : %s", sErr)
			}
			break
		}
	} else {
		if err := startCmd(d, cmd); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func initDriver(d *basic.Driver) error {
	dbusConfig, ok := d.Config.CustomOptions["dbus"]
	if !ok {
		return fmt.Errorf("Failed to find dbus params in driver custom options")
	}
	if dbusConfig.(string) == "" {
		// start needed dbus like (eval `dbus-launch --sh-syntax`) and get env variable
		result, err := exec.Command("dbus-launch", "--sh-syntax").CombinedOutput() //DBUS_SESSION_BUS_ADDRESS='unix:abstract=/tmp/dbus-JHGXLpeJ6A,guid=25ab632502ebccd43cd403bc58388fab';\n ...
		if err != nil {
			panic(err)
		}
		env := string(result)
		logrus.Debugf("dbus-launch --sh-syntax -> \n%s", env)
		reDBus := regexp.MustCompile("DBUS_SESSION_BUS_ADDRESS='(.*?)';")
		//rePID := regexp.MustCompile("DBUS_SESSION_BUS_PID=(.*?);")
		matchDBuse := reDBus.FindStringSubmatch(env)
		//matchPID := rePID.FindStringSubmatch(env)
		d.Config.CustomOptions["dbus"] = matchDBuse[1]
		//TODO plan to kill this add closing ?
	}

	d.Config.CustomOptions["env"] = []string{fmt.Sprintf("DBUS_SESSION_BUS_ADDRESS=%s", d.Config.CustomOptions["dbus"])}
	return startFuseDeamon(d)
}

func startFuseDeamon(d *GVfsDriver) error {
	//TODO check needed gvfsd + gvfsd-ftp Maybe already on dbus ?
	// Normaly gvfsd-fuse block such so this like crash but global ?

	fi, err := os.Lstat(d.Config.Root)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(d.Config.Root, 0700); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	if fi != nil && !fi.IsDir() {
		return fmt.Errorf("%v already exist and it's not a directory", d.Config.Root)
	}

	err = startCmd(d, "/usr/lib/gvfs/gvfsd --no-fuse") //Start global deamon
	if err != nil {
		return err
	}

	err = startCmd(d, fmt.Sprintf("/usr/lib/gvfs/gvfsd-fuse %s -f -o %s", d.Config.Root, d.Config.CustomOptions["fuseOpts"])) //Start ftp handler
	return err
}

// start deamon in context of this gvfs drive with custome env
func startCmd(d *GVfsDriver, cmd string) error {
	logrus.Debugln(d.Config.CustomOptions["env"].([]string), cmd)
	return setEnv(cmd, d.Config.CustomOptions["env"].([]string)).Start()
}
