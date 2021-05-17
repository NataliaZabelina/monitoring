package command

import (
	"testing"

	"github.com/NataliaZabelina/monitoring/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestGetLoadAvg(t *testing.T) {
	db := &storage.Db{}
	db.Init()
	result, err := GetLoadAvg(db)
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Positive(t, db.System_table.GetAverage(5).Load_avg)
}

func TestGetCpu(t *testing.T) {
	db := &storage.Db{}
	db.Init()
	result, err := GetCpu(db)
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.NotEmpty(t, db.CPU_table.GetAverage(5).User_mode)
	require.NotEmpty(t, db.CPU_table.GetAverage(5).System_mode)
	require.NotEmpty(t, db.CPU_table.GetAverage(5).Idle)
}
