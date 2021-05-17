package parser

import (
	"runtime"

	"github.com/NataliaZabelina/monitoring/internal/storage"
)

var os = runtime.GOOS

func GetParser(command string) func(*storage.Db, string) error {
	switch command {
	case "load_avg":
		return getLoadAvgParser()
	case "cpu":
		return getCPUParser()
	// case "disk_io":
	// 	return getDiskIOParser()
	default:
		return getLoadAvgParser()
	}
}

func getLoadAvgParser() func(*storage.Db, string) error {
	parserDarwin := ParseSystemLoadDarwin
	parserLinux := ParseSystemLoadLinux
	return chooseParser(parserDarwin, parserLinux)
}

func getCPUParser() func(*storage.Db, string) error {
	parserDarwin := ParseCPULoadDarwin
	parserLinux := ParseCpuLoadLinux
	return chooseParser(parserDarwin, parserLinux)
}

// func getDiskIOParser() func(*storage.Db, string) error {
// 	parserDarwin := "iostat -dC"
// 	parserLinux := "iostat -d -k"
// 	return chooseParser(parserDarwin, parserLinux)
// }

// func getDiskFSParser() func(*storage.Db, string) error {
// 	parserDarwin := "df -m"
// 	parserLinux := "df -i"
// 	return chooseParser(parserDarwin, parserLinux)
// }

// func getTopTalkersParser() func(*storage.Db, string) error {
// 	parserDarwin := ""      //todo
// 	parserLinux := "ss -ta" //todo
// 	return chooseParser(parserDarwin, parserLinux)
// }

// func getNetStatParser() func(*storage.Db, string) error {
// 	parserDarwin := ""                   // todo
// 	parserLinux := "sudo netstat -lntup" //todo
// 	return chooseParser(parserDarwin, parserLinux)
// }

func chooseParser(parserDarwin func(*storage.Db, string) error, 
	parserLinux func(*storage.Db, string) error) func(*storage.Db, string) error {
	switch os {
	case "darwin":
		return parserDarwin
	case "linux":
		return parserLinux
	default:
		return parserLinux
	}
}
