package drivers

import "testing"

func TestisFileRoot(t *testing.T) {
	is, _ := isFile("/") //Is Folder so shouldn't be ok
	if is {
		t.Error("Expected false, got ", is)
	}
}

func TesturlToDriverFTP(t *testing.T) {
	d, err := urlToDriver("ftp://host") //Is Folder so shouldn't be ok
	if err != nil || d.id() != FTP {
		t.Error("Expected a FTP driver, got ", d, err)
	}
}
func TesturlToDriverSSH(t *testing.T) {
	d, err := urlToDriver("ssh://host")
	if err != nil || d.id() != SSH {
		t.Error("Expected a SSH driver, got ", d, err)
	}
}
func TesturlToDriverSFTP(t *testing.T) {
	d, err := urlToDriver("sftp://host")
	if err != nil || d.id() != SSH {
		t.Error("Expected a SSH driver, got ", d, err)
	}
}
func TesturlToDriverSMB(t *testing.T) {
	d, err := urlToDriver("smb://host")
	if err != nil || d.id() != SMB {
		t.Error("Expected a SMB driver, got ", d, err)
	}
}
func TesturlToDriverDav(t *testing.T) {
	d, err := urlToDriver("dav://host")
	if err != nil || d.id() != DAV {
		t.Error("Expected a SMB driver, got ", d, err)
	}
}
func TesturlToDriverDavs(t *testing.T) {
	d, err := urlToDriver("davs://host")
	if err != nil || d.id() != DAV {
		t.Error("Expected a SMB driver, got ", d, err)
	}
}
func TesturlToDriverError(t *testing.T) {
	d, err := urlToDriver("host")
	if err == nil {
		t.Error("Expected a error, got ", d, err)
	}
}
