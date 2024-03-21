package env

type cpuArch uint8

const (
	bit32 cpuArch = iota
	bit64
)

type distro string

const (
	debian distro = "debian"
	ubuntu distro = "ubuntu"
)

var supportedOS = map[string]distro{
	"debian": debian,
	"ubuntu": ubuntu,
}

var (
	osDistro distro
	// osVer    string
)
