package command

import "runtime"

var os = runtime.GOOS

func GetCommand(command string) string {
	switch command {
	case "load_avg":
		return getLoadAvgCommand()
	case "cpu":
		return getCPUCommand()
	case "disk_io":
		return getDiskIOCommand()
	default:
		return getLoadAvgCommand()
	}
}

func getLoadAvgCommand() string {
	cmdDarwin := "top | head -3 | tail -1 | cut -d\":\" -f2"
	cmdLinux := "top -b -n1"
	return chooseCmd(cmdDarwin, cmdLinux)
}

func getCPUCommand() string {
	cmdDarwin := "top | head -4 | tail -1 | cut -d\":\" -f2"
	cmdLinux := "top -b -n1"
	return chooseCmd(cmdDarwin, cmdLinux)
}

func getDiskIOCommand() string {
	cmdDarwin := "iostat -dC"
	cmdLinux := "iostat -d -k"
	return chooseCmd(cmdDarwin, cmdLinux)
}

func getDiskFSCommand() string {
	cmdDarwin := "df -m"
	cmdLinux := "df -i"
	return chooseCmd(cmdDarwin, cmdLinux)
}

func getTopTalkersCommand() string {
	cmdDarwin := ""      //todo
	cmdLinux := "ss -ta" //todo
	return chooseCmd(cmdDarwin, cmdLinux)
}

func getNetStatCommand() string {
	cmdDarwin := ""                   // todo
	cmdLinux := "sudo netstat -lntup" //todo
	return chooseCmd(cmdDarwin, cmdLinux)
}

func chooseCmd(cmdDarwin string, cmdLinux string) string {
	switch os {
	case "darwin":
		return cmdDarwin
	case "linux":
		return cmdLinux
	default:
		return cmdLinux
	}
}
