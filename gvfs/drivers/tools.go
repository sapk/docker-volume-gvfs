package drivers

import (
	"fmt"
	"net/url"
	"os"
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

// Check if path exist and is not a folder
func isFile(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		if f.IsDir() {
			return false, fmt.Errorf("File is a folder not a binary")
		}
		return true, nil
	}
	return false, err
}

func urlToDriver(urlStr string) (gvfsVolumeDriver, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "ftp":
	case "ftps":
		return FTPVolumeDriver{url: u}, nil
	case "ssh":
	case "sftp":
		return SSHVolumeDriver{url: u}, nil
	case "smb":
		return SMBVolumeDriver{url: u}, nil
	case "dav":
	case "davs":
		return DavVolumeDriver{url: u}, nil
	}
	return nil, fmt.Errorf("%v is not matching any known driver", urlStr)
}
