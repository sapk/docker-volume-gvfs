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
)

var driverTypes = []string{
	"ftp",
	"ssh",
	"smb",
}

func (dt DriverType) String() string {
	return driverTypes[dt]
}
