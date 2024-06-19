package collector

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestAccountMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("fixtures/squeue/account.txt")
	data, _ := io.ReadAll(file)
	accounts := ParseAccountMetrics(data)

	assert.Equal(t, 35.0, accounts["ampere"].pending, "Miscount of pending account jobs")
	assert.Equal(t, 152.0, accounts["ampere"].pendingCpus, "Miscount of cpusPending account jobs")
	assert.Equal(t, 30.0, accounts["ampere"].running, "Miscount of running account jobs")
	assert.Equal(t, 269.0, accounts["ampere"].runningCpus, "Miscount of runningCpus account jobs")
	assert.Equal(t, 0.0, accounts["ampere"].suspended, "Miscount of suspended account jobs")
}
