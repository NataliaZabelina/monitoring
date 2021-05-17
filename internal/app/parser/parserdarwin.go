package parser

import (
	"strconv"
	"strings"

	l "github.com/NataliaZabelina/monitoring/internal/logger"
	"github.com/NataliaZabelina/monitoring/internal/storage"
	"github.com/NataliaZabelina/monitoring/internal/storage/schema"
)

func ParseSystemLoadDarwin(db *storage.Db, txt string) error {
	arr := strings.Split(txt, ",")
	if len(arr) > 0 {
		val, err := strconv.ParseFloat(arr[0], 32)
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

func ParseCPULoadDarwin(db *storage.Db, txt string) error {
	arr := strings.Split(txt, ",")
	if len(arr) > 0 {
		um_val, err := TrimLeftAndGetFLoat(arr[0], "%")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu user mode stat: %w", err)
			return err
		}

		sm_val, err := TrimLeftAndGetFLoat(arr[1], "%")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu system mode stat: %w", err)
			return err
		}

		idle_val, err := TrimLeftAndGetFLoat(arr[2], "%")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu idle stat: %w", err)
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
