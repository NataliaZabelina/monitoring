package parser

import (
	"testing"

	"github.com/NataliaZabelina/monitoring/internal/storage"
	"github.com/NataliaZabelina/monitoring/internal/storage/schema"
	"github.com/stretchr/testify/require"
)

var (
	outputIostat = `Linux 5.4.0-1042-gcp (project-d8034bc9-c018-4c7b-815a-6d1689973512)     05/17/21        _x86_64_        (4 CPU)
    
    Device             tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
    loop0             0.06         0.08         0.00         0.00       1884          0          0
    loop1             0.00         0.04         0.00         0.00       1077          0          0
    loop2             0.00         0.04         0.01         0.00       1093          0          0
    sda              36.58        44.34       657.54       645.55    1113005   16505437   16204630`

	expectedIostat = schema.DiskDto{
		Disk: map[string]schema.Disk{
			"loop0": {Param1: 0.06, Param2: 0.08},
			"loop1": {Param1: 0.00, Param2: 0.04},
			"loop2": {Param1: 0.00, Param2: 0.05},
			"sda":   {Param1: 36.58, Param2: 701.88},
		},
	}
)

func TestDiskIOParserLinux(t *testing.T) {
	db := &storage.DB{}
	db.Init()
	err := ParseDiskIOLinux(db, outputIostat)
	require.Nil(t, err)
	current := db.DiskTable.GetAverage(1)
	require.InDelta(t, expectedIostat.Disk["loop0"].Param1, current.Disk["loop0"].Param1, 0.00001)
	require.InDelta(t, expectedIostat.Disk["loop0"].Param2, current.Disk["loop0"].Param2, 0.00001)
	require.InDelta(t, expectedIostat.Disk["loop1"].Param1, current.Disk["loop1"].Param1, 0.00001)
	require.InDelta(t, expectedIostat.Disk["loop1"].Param2, current.Disk["loop1"].Param2, 0.00001)
	require.InDelta(t, expectedIostat.Disk["loop2"].Param1, current.Disk["loop2"].Param1, 0.00001)
	require.InDelta(t, expectedIostat.Disk["loop2"].Param2, current.Disk["loop2"].Param2, 0.00001)
	require.InDelta(t, expectedIostat.Disk["sda"].Param1, current.Disk["sda"].Param1, 0.00001)
	require.InDelta(t, expectedIostat.Disk["sda"].Param2, current.Disk["sda"].Param2, 0.00001)
}

var (
	outputTop = `top - 02:17:38 up 18 days,  1:54,  0 users,  load average: 0.23, 0.17, 0.12
	Tasks:   2 total,   1 running,   1 sleeping,   0 stopped,   0 zombie
    %Cpu(s):  0.0 us,  1.1 sy,  0.0 ni, 98.9 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
    MiB Mem :   1989.2 total,     82.9 free,    491.8 used,   1414.5 buff/cache
    MiB Swap:   1024.0 total,    987.1 free,     36.9 used.   1399.7 avail Mem 
	
      PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
          1 root      20   0    3992   2784   2336 S   0.0   0.1   0:00.12 bash
     7537 root      20   0    7864   2996   2636 R   0.0   0.1   0:00.00 top`

	expectedCPU = schema.CPULoadDto{
		UserMode:   0.0,
		SystemMode: 1.1,
		Idle:       98.9,
	}

	expectedLoadAvg = schema.SystemLoadDto{
		LoadAvg: 0.23,
	}
)

func TestLoadAvgParserLinux(t *testing.T) {
	db := &storage.DB{}
	db.Init()
	err := ParseSystemLoadLinux(db, outputTop)
	require.Nil(t, err)
	current := db.SystemTable.GetAverage(1)
	require.InDelta(t, expectedLoadAvg.LoadAvg, current.LoadAvg, 0.00001)
}

func TestCPUParserLinux(t *testing.T) {
	db := &storage.DB{}
	db.Init()
	err := ParseCPULoadLinux(db, outputTop)
	require.Nil(t, err)
	current := db.CPUTable.GetAverage(1)
	require.InDelta(t, expectedCPU.UserMode, current.UserMode, 0.00001)
	require.InDelta(t, expectedCPU.SystemMode, current.SystemMode, 0.00001)
	require.InDelta(t, expectedCPU.Idle, current.Idle, 0.00001)
}
