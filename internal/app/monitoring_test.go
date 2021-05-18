package monitoring

import (
	"testing"

	"github.com/NataliaZabelina/monitoring/internal/app/command"
	"github.com/NataliaZabelina/monitoring/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestGetLoadAvg(t *testing.T) {
	db := &storage.DB{}
	db.Init()
	result, err := command.GetLoadAvg(db)
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Positive(t, db.SystemTable.GetAverage(5).LoadAvg)
}

func TestGetCPU(t *testing.T) {
	db := &storage.DB{}
	db.Init()
	result, err := command.GetCPU(db)
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.GreaterOrEqual(t, float64(db.CPUTable.GetAverage(5).UserMode), 0.0)
	require.GreaterOrEqual(t, float64(db.CPUTable.GetAverage(5).SystemMode), 0.0)
	require.GreaterOrEqual(t, float64(db.CPUTable.GetAverage(5).Idle), 0.0)
}

func TestGetDiskIO(t *testing.T) {
	db := &storage.DB{}
	db.Init()
	result, err := command.GetDiskIO(db)
	require.Nil(t, err)
	require.NotEmpty(t, result)
	disk := db.DiskTable.GetAverage(5).Disk
	for _, v := range disk {
		require.GreaterOrEqual(t, float64(v.Param1), 0.0)
		require.GreaterOrEqual(t, float64(v.Param2), 0.0)
	}
}
