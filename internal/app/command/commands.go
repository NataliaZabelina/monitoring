package command

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/NataliaZabelina/monitoring/internal/app/parser"
	l "github.com/NataliaZabelina/monitoring/internal/logger"
	"github.com/NataliaZabelina/monitoring/internal/storage"
)

func GetLoadAvg(db *storage.DB) (string, error) {
	return runCommand(db, "load_avg")
}

func GetCPU(db *storage.DB) (string, error) {
	return runCommand(db, "cpu")
}

func GetDiskIO(db *storage.DB) (string, error) {
	return runCommand(db, "disk_io")
}

func GetDiskFS(db *storage.DB) (string, error) {
	return runCommand(db, "disk_fs")
}

func GetTopTalkers(db *storage.DB) (string, error) {
	return runCommand(db, "top_talkers")
}

func GetNetStat(db *storage.DB) (string, error) {
	return runCommand(db, "net_stat")
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

func runCommand(db *storage.DB, command string) (string, error) {
	out, err := execute("/bin/bash", "-c", GetCommand(command))
	if err != nil {
		return "", err
	}
	result := strings.Trim(out, " ")
	err = parser.GetParser(command)(db, result)
	if err != nil {
		return result, err
	}
	return result, nil
}
