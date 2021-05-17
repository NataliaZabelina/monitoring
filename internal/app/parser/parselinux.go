package parser

import (
	"strconv"
	"strings"

	l "github.com/NataliaZabelina/monitoring/internal/logger"
	"github.com/NataliaZabelina/monitoring/internal/storage"
	"github.com/NataliaZabelina/monitoring/internal/storage/schema"
)

func ParseSystemLoadLinux(db *storage.Db, txt string) error {
	loads := splitInput(txt, ":")

	if len(loads) > 0 {
		txt = strings.Replace(loads[0], ",", ".", 1)
		val, err := strconv.ParseFloat(txt, 32)
		if err != nil {
			l.Logger.Errorf("Can't parse load average stat: %w", err)
			return err
		}
		db.System_table.AddEntity(schema.SystemLoadDto{
			Load_avg: float32(val),
		})
	}
	return nil
}

func ParseCpuLoadLinux(db *storage.Db, txt string) error {
	cpu := splitInput(txt, ":")
	if len(cpu) > 0 {
		um_str := strings.Replace(cpu[0], ",", ".", 1)
		um_val, err := TrimLeftAndGetFLoat(um_str, " us")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu user mode stat: %w", err)
			return err
		}

		sm_str := strings.Replace(cpu[1], ",", ".", 1)
		sm_val, err := TrimLeftAndGetFLoat(sm_str, " sy")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu system mode stat: %w", err)
			return err
		}

		idle_str := strings.Replace(cpu[3], ",", ".", 1)
		idle_val, err := TrimLeftAndGetFLoat(idle_str, " id")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu user mode stat: %w", err)
			return err
		}

		db.CPU_table.AddEntity(schema.CpuLoadDto{
			User_mode:   um_val,
			System_mode: sm_val,
			Idle:        idle_val,
		})
	}
	return nil
}

func splitInput(str string, trimmer string) []string {
	last_index := strings.LastIndex(str, trimmer)
	str = strings.TrimLeft(str[last_index+1:], " ")
	arr := strings.Split(str, ", ")
	return arr
}

func TrimLeftAndGetFLoat(str string, trimmer string) (float32, error) {
	last_index := strings.LastIndex(str, trimmer)
	str_val := strings.TrimLeft(str[:last_index], " ")
	val, err := strconv.ParseFloat(str_val, 32)
	return float32(val), err
}
