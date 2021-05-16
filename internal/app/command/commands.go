package command

import (
	"bytes"
	"os/exec"
	"strings"

	l "github.com/NataliaZabelina/monitoring/internal/logger"
)

func GetLoadAvg() string {
	return runCommand("load_avg")
}

func GetCpu() string {
	return runCommand("cpu")
}

func GetDiskIO() string {
	return runCommand("disk_io")
}

func GetDiskFS() string {
	return runCommand("disk_fs")
}

func GetTopTalkers() string {
	return runCommand("top_talkers")
}

func GetNetStat() string {
	return runCommand("net_stat")
}

func execute(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		l.Logger.Errorf("Error when execute command %s: $v", name, err)
		return stdout.String(), err
	}
	return stdout.String(), nil
}

func runCommand(command string) string {
	out, err := execute("/bin/bash", "-c", GetCommand(command))
	if err != nil {
		return ""
	}
	result := strings.Trim(out, " ")
	return result
}
