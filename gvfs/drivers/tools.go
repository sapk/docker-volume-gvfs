package drivers

import (
	"fmt"
	"net/url"
	"os/exec"
)

func setEnv(cmd string, env []string) *exec.Cmd {
	c := exec.Command("sh", "-c", cmd)
	//c := exec.Command(strings.Split(cmd, " ")) //TODO better
	c.Env = env
	return c
}

func getDriver(urlStr string) (*gvfsVolumeDriver, string, error) {
	d, err := urlToDriver(urlStr)
	if err != nil {
		return nil, "", err
	}
	if !d.isAvailable() {
		return nil, "", fmt.Errorf("%s driver is not available on this host", d.id())
	}
	m, err := d.mountpoint()
	if err != nil {
		return nil, "", err
	}
	return &d, m, nil
}

func urlToDriver(urlStr string) (gvfsVolumeDriver, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "ftp":
		return FTPVolumeDriver{url: u}, nil
	case "ssh":
		return SSHVolumeDriver{url: u}, nil
	case "smb":
		return SMBVolumeDriver{url: u}, nil
	default:
		return nil, fmt.Errorf("%v is not matching any known driver", urlStr)
	}
}
