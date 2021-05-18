package parser

import (
	"bufio"
	"strconv"
	"strings"

	l "github.com/NataliaZabelina/monitoring/internal/logger"
	"github.com/NataliaZabelina/monitoring/internal/storage"
	"github.com/NataliaZabelina/monitoring/internal/storage/schema"
)

func ParseSystemLoadLinux(db *storage.DB, txt string) error {
	loads := splitInput(txt, "load average:")

	if len(loads) > 0 {
		txt = strings.Replace(loads[0], ",", ".", 1)
		val, err := strconv.ParseFloat(txt, 32)
		if err != nil {
			l.Logger.Errorf("Can't parse load average stat: %w", err)
			return err
		}
		db.SystemTable.AddEntity(schema.SystemLoadDto{
			LoadAvg: float32(val),
		})
	}
	return nil
}

func ParseCPULoadLinux(db *storage.DB, txt string) error {
	cpu := splitInput(txt, "%Cpu(s):")
	if len(cpu) > 0 {
		umStr := strings.Replace(cpu[0], ",", ".", 1)
		umVal, err := TrimLeftAndGetFLoat(umStr, " us")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu user mode stat: %w", err)
			return err
		}

		smStr := strings.Replace(cpu[1], ",", ".", 1)
		smVal, err := TrimLeftAndGetFLoat(smStr, " sy")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu system mode stat: %w", err)
			return err
		}

		idleStr := strings.Replace(cpu[3], ",", ".", 1)
		idleVal, err := TrimLeftAndGetFLoat(idleStr, " id")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu user mode stat: %w", err)
			return err
		}

		db.CPUTable.AddEntity(schema.CPULoadDto{
			UserMode:   umVal,
			SystemMode: smVal,
			Idle:       idleVal,
		})
	}
	return nil
}

func ParseDiskIOLinux(db *storage.DB, txt string) error {
	output := bufio.NewScanner(strings.NewReader(txt))
	result := map[string]schema.Disk{}
	for output.Scan() {
		token := strings.TrimLeft(output.Text(), " ")
		if !strings.HasPrefix(token, "Device") {
			continue
		}
		for output.Scan() {
			s := strings.ReplaceAll(output.Text(), ",", ".")
			data := strings.Fields(s)
			if len(data) < 6 {
				continue
			}

			tps, err := strconv.ParseFloat(data[1], 32)
			if err != nil {
				l.Logger.Errorf("Can't parse tps stat: %w", err)
				return err
			}
			kbReadS, err := strconv.ParseFloat(data[2], 32)
			if err != nil {
				l.Logger.Errorf("Can't parse kbReadS stat: %w", err)
				return err
			}
			kbWriteS, err := strconv.ParseFloat(data[3], 32)
			if err != nil {
				l.Logger.Errorf("Can't parse kbWriteS stat: %w", err)
				return err
			}

			diskIO := schema.Disk{
				Param1: float32(tps),
				Param2: float32(kbReadS + kbWriteS),
			}

			result[data[0]] = diskIO
		}
	}
	db.DiskTable.AddEntity(schema.DiskDto{
		Disk: result,
	})
	return nil
}

func ParseDiskFSLinux(db *storage.DB, txt string) error {
	// scanner := bufio.NewScanner(strings.NewReader(txt))

	// for scanner.Scan() {
	// 	if !strings.HasPrefix(scanner.Text(), "Filesystem") {
	// 		continue
	// 	}

	// 	for scanner.Scan() {
	// 		s := strings.ReplaceAll(scanner.Text(), ",", ".")
	// 		data := strings.Fields(s)
	// 		if len(data) < 6 {
	// 			continue
	// 		}
	// 		used := strings.TrimLeft(data[4], "%")
	// 		usedProc, err := strconv.Atoi(used)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		diskFS := schema.Disk{
	// 			Param1: float32(usedProc),
	// 			Param2: float32(),
	// 		}

	// 		result := map[string]schema.Disk{}
	// 		result[data[0]] = diskFS

	// 		db.DiskTable.AddEntity(schema.DiskDto{
	// 			Disk: result,
	// 		})
	// 	}
	// }

	return nil
}

func splitInput(str string, trimmer string) []string {
	lastIndex := strings.LastIndex(str, trimmer) + len(trimmer) - 1
	str = strings.TrimLeft(str[lastIndex+1:], " ")
	arr := strings.Split(str, ",")
	return arr
}

func TrimLeftAndGetFLoat(str string, trimmer string) (float32, error) {
	lastIndex := strings.LastIndex(str, trimmer)
	strVal := strings.TrimLeft(str[:lastIndex], " ")
	val, err := strconv.ParseFloat(strVal, 32)
	return float32(val), err
}
