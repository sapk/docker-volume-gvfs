package drivers

//DriverType id of type of driver
type DriverType int

const (
	//FTP id of ftp driver type
	FTP DriverType = iota
	//SSH id of ssh driver type
	SSH
	//SMB id of smb driver type
	SMB
	//DAV id of dav driver type
	DAV
	//NFS id of nfs driver type
	NFS
)

var driverTypes = []string{
	"ftp",
	"ssh",
	"smb",
	"dav",
	"nfs",
}

func (dt DriverType) String() string {
	return driverTypes[dt]
}
