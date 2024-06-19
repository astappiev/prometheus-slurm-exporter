package collector

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestParsePartitionsMetrics(t *testing.T) {
	// Read the input data from a file
	sinfoFile, _ := os.Open("fixtures/sinfo/partition.txt")
	runningFile, _ := os.Open("fixtures/squeue/partition_running.txt")
	pendingFile, _ := os.Open("fixtures/squeue/partition_pending.txt")
	sinfoData, _ := io.ReadAll(sinfoFile)
	runningData, _ := io.ReadAll(runningFile)
	pendingData, _ := io.ReadAll(pendingFile)
	partitionMetrics := ParsePartitionMetrics(sinfoData, runningData, pendingData)

	assert.Equal(t, 273.0, partitionMetrics["ampere"].cpu.alloc, "Miscount of allocated CPUs")
	assert.Equal(t, 279.0, partitionMetrics["ampere"].cpu.idle, "Miscount of idle CPUs")
	assert.Equal(t, 480.0, partitionMetrics["ampere"].cpu.other, "Miscount of other CPUs")
	assert.Equal(t, 1032.0, partitionMetrics["ampere"].cpu.total, "Miscount of total CPUs")
	assert.Equal(t, 16.0, partitionMetrics["ampere"].jobsRunning, "Miscount of running jobs")
	assert.Equal(t, 30.0, partitionMetrics["ampere"].jobsPending, "Miscount of pending jobs")
}
