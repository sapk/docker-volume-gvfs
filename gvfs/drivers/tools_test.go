package drivers

import "testing"

func TestIsFileRoot(t *testing.T) {
	is, err := isFile("/") //Is Folder so shouldn't be ok
	if is {
		t.Error("Expected false, got ", is, err)
	}
}
func TestIsFileNotExist(t *testing.T) {
	is, err := isFile("/dshdskjdkjs/sdjsjdsjks") //ls should be nok on all host
	if is {
		t.Error("Expected false, got ", is, err)
	}
}
func TestIsFileLS(t *testing.T) {
	is, err := isFile("/usr/bin/ls") //ls should be ok on all host
	if !is || err != nil {
		t.Error("Expected true without error, got ", is, err)
	}
}

func TestUrlToDriverFTP(t *testing.T) {
	d, err := urlToDriver("ftp://host") //Is Folder so shouldn't be ok
	if err != nil || d.id() != FTP {
		t.Error("Expected a FTP driver, got ", d, err)
	}
}
func TestUrlToDriverSSH(t *testing.T) {
	d, err := urlToDriver("ssh://host")
	if err != nil || d.id() != SSH {
		t.Error("Expected a SSH driver, got ", d, err)
	}
}
func TestUrlToDriverSFTP(t *testing.T) {
	d, err := urlToDriver("sftp://host")
	if err != nil || d.id() != SSH {
		t.Error("Expected a SSH driver, got ", d, err)
	}
}
func TestUrlToDriverSMB(t *testing.T) {
	d, err := urlToDriver("smb://host")
	if err != nil || d.id() != SMB {
		t.Error("Expected a SMB driver, got ", d, err)
	}
}
func TestUrlToDriverDav(t *testing.T) {
	d, err := urlToDriver("dav://host")
	if err != nil || d.id() != DAV {
		t.Error("Expected a DAV driver, got ", d, err)
	}
}
func TestUrlToDriverDavs(t *testing.T) {
	d, err := urlToDriver("davs://host")
	if err != nil || d.id() != DAV {
		t.Error("Expected a DAV driver, got ", d, err)
	}
}
func TestUrlToDriverNFS(t *testing.T) {
	d, err := urlToDriver("nfs://host")
	if err != nil || d.id() != NFS {
		t.Error("Expected a NFS driver, got ", d, err)
	}
}
func TestUrlToDriverUnkownError(t *testing.T) {
	d, err := urlToDriver("ltp://host")
	if err == nil {
		t.Error("Expected a error, got ", d, err)
	}
}
func TestUrlToDriverURLError(t *testing.T) {
	d, err := urlToDriver("host")
	if err == nil {
		t.Error("Expected a error, got ", d, err)
	}
}
