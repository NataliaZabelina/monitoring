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
			"loop0": schema.Disk{
				Param1: 0.06,
				Param2: 0.08,
			},
			"loop1": schema.Disk{
				Param1: 0.0,
				Param2: 0.04,
			},
			"loop2": schema.Disk{
				Param1: 0.00,
				Param2: 0.05,
			},
			"sda": schema.Disk{
				Param1: 36.58,
				Param2: 701.88,
			},
		},
	}
)

func TestDiskIOParser(t *testing.T) {
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
