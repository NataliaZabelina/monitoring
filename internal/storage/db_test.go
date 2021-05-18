package storage

import (
	"testing"

	"github.com/NataliaZabelina/monitoring/internal/storage/schema"
	"github.com/stretchr/testify/require"
)

func TestSystemGetAverage(t *testing.T) {
	db := &DB{}
	db.Init()
	db.SystemTable.AddEntity(schema.SystemLoadDto{
		LoadAvg: 3.8,
	})
	db.SystemTable.AddEntity(schema.SystemLoadDto{
		LoadAvg: 1.0,
	})
	out := db.SystemTable.GetAverage(1)

	require.InDelta(t, 2.4, out.LoadAvg, 0.00001)
}

func TestCPUGetAverage(t *testing.T) {
	db := &DB{}
	db.Init()
	db.CPUTable.AddEntity(schema.CPULoadDto{
		UserMode:   1.4,
		SystemMode: 0.9,
		Idle:       87.0,
	})
	db.CPUTable.AddEntity(schema.CPULoadDto{
		UserMode:   1.6,
		SystemMode: 0.3,
		Idle:       45.3,
	})
	out := db.CPUTable.GetAverage(1)

	require.InDelta(t, 1.5, out.UserMode, 0.00001)
	require.InDelta(t, 0.6, out.SystemMode, 0.00001)
	require.InDelta(t, 66.15, out.Idle, 0.00001)
}

func TestDiskGetAverage(t *testing.T) {
	db := &DB{}
	db.Init()

	disks1 := map[string]schema.Disk{
		"loop0": {Param1: 1.60, Param2: 0.08},
		"loop1": {Param1: 0.00, Param2: 4.74},
		"loop2": {Param1: 0.00, Param2: 0.00},
		"sda":   {Param1: 66.51, Param2: 68.52},
	}

	disks2 := map[string]schema.Disk{
		"loop0": {Param1: 1.13, Param2: 0.5},
		"loop1": {Param1: 0.30, Param2: 0.74},
		"loop2": {Param1: 0.00, Param2: 0.60},
		"sda":   {Param1: 68.50, Param2: 70.32},
	}

	db.DiskTable.AddEntity(schema.DiskDto{
		Disk: disks1,
	})

	db.DiskTable.AddEntity(schema.DiskDto{
		Disk: disks2,
	})

	out := db.DiskTable.GetAverage(1)

	for k, v := range out.Disk {
		switch k {
		case "loop0":
			require.InDelta(t, 1.365, v.Param1, 0.00001)
			require.InDelta(t, 0.29, v.Param2, 0.00001)
		case "loop1":
			require.InDelta(t, 0.15, v.Param1, 0.00001)
			require.InDelta(t, 2.74, v.Param2, 0.00001)
		case "loop2":
			require.InDelta(t, 0.00, v.Param1, 0.00001)
			require.InDelta(t, 0.30, v.Param2, 0.00001)
		case "sda":
			require.InDelta(t, 67.505, v.Param1, 0.00001)
			require.InDelta(t, 69.42, v.Param2, 0.00001)
		}
	}
}
