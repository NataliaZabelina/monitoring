package parser

import (
	"strconv"
	"strings"

	l "github.com/NataliaZabelina/monitoring/internal/logger"
	"github.com/NataliaZabelina/monitoring/internal/storage"
	"github.com/NataliaZabelina/monitoring/internal/storage/schema"
)

func ParseSystemLoadDarwin(db *storage.DB, txt string) error {
	arr := strings.Split(txt, ",")
	if len(arr) > 0 {
		val, err := strconv.ParseFloat(arr[0], 32)
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

func ParseCPULoadDarwin(db *storage.DB, txt string) error {
	arr := strings.Split(txt, ",")
	if len(arr) > 0 {
		umVal, err := TrimLeftAndGetFLoat(arr[0], "%")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu user mode stat: %w", err)
			return err
		}

		smVal, err := TrimLeftAndGetFLoat(arr[1], "%")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu system mode stat: %w", err)
			return err
		}

		idleVal, err := TrimLeftAndGetFLoat(arr[2], "%")
		if err != nil {
			l.Logger.Errorf("Can't parse cpu idle stat: %w", err)
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

func ParseDiskIODarwin(db *storage.DB, txt string) error {
	return nil
}
