package parser

import (
	"runtime"

	"github.com/NataliaZabelina/monitoring/internal/storage"
)

var os = runtime.GOOS

func GetParser(command string) func(*storage.DB, string) error {
	switch command {
	case "load_avg":
		return getLoadAvgParser()
	case "cpu":
		return getCPUParser()
	case "disk_io":
		return getDiskIOParser()
	case "disk_fs":
		return getDiskFSParser()
	default:
		return getLoadAvgParser()
	}
}

func getLoadAvgParser() func(*storage.DB, string) error {
	parserDarwin := ParseSystemLoadDarwin
	parserLinux := ParseSystemLoadLinux
	return chooseParser(parserDarwin, parserLinux)
}

func getCPUParser() func(*storage.DB, string) error {
	parserDarwin := ParseCPULoadDarwin
	parserLinux := ParseCPULoadLinux
	return chooseParser(parserDarwin, parserLinux)
}

func getDiskIOParser() func(*storage.DB, string) error {
	parserDarwin := ParseDiskIODarwin
	parserLinux := ParseDiskIOLinux
	return chooseParser(parserDarwin, parserLinux)
}

func getDiskFSParser() func(*storage.DB, string) error {
	parserDarwin := ParseDiskFSDarwin
	parserLinux := ParseDiskFSLinux
	return chooseParser(parserDarwin, parserLinux)
}

func chooseParser(parserDarwin func(*storage.DB, string) error,
	parserLinux func(*storage.DB, string) error) func(*storage.DB, string) error {
	switch os {
	case "darwin":
		return parserDarwin
	case "linux":
		return parserLinux
	default:
		return parserLinux
	}
}
