package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecuteCPU(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{
			name:    "Get LoadAvg output",
			command: "load_avg",
		},
		{
			name:    "Get CPU output",
			command: "cpu",
		},
		{
			name:    "Get DiskIO output",
			command: "disk_io",
		},
	}

	for _, testCase := range testCases {
		c := testCase.command
		t.Run(testCase.name, func(t *testing.T) {
			cmd := GetCommand(c)
			out, err := execute("/bin/bash", "-c", cmd)
			require.Nil(t, err)
			require.NotNil(t, out)
		})
	}
}
