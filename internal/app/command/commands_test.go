package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLoadAvg(t *testing.T) {
	result := GetLoadAvg()
	require.NotEmpty(t, result)
}

func TestGetCpu(t *testing.T) {
	result := GetCpu()
	require.NotEmpty(t, result)
}