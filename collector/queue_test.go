package collector

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestParseQueueMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("fixtures/squeue/queue.txt")
	data, _ := io.ReadAll(file)
	queueMetrics := ParseQueueMetrics(data)

	assert.Equal(t, 32.0, queueMetrics.pending, "Miscount of pending jobs")
	assert.Equal(t, 0.0, queueMetrics.pendingDep, "Miscount of pendingDep jobs")
	assert.Equal(t, 27.0, queueMetrics.running, "Miscount of running jobs")
	assert.Equal(t, 0.0, queueMetrics.suspended, "Miscount of suspended jobs")
	assert.Equal(t, 1.0, queueMetrics.cancelled, "Miscount of cancelled jobs")
	assert.Equal(t, 1.0, queueMetrics.completing, "Miscount of completing jobs")
	assert.Equal(t, 1.0, queueMetrics.completed, "Miscount of completed jobs")
	assert.Equal(t, 0.0, queueMetrics.configuring, "Miscount of configuring jobs")
	assert.Equal(t, 0.0, queueMetrics.failed, "Miscount of failed jobs")
	assert.Equal(t, 0.0, queueMetrics.timeout, "Miscount of timeout jobs")
	assert.Equal(t, 0.0, queueMetrics.preempted, "Miscount of preempted jobs")
	assert.Equal(t, 0.0, queueMetrics.nodeFail, "Miscount of nodeFail jobs")
	assert.Equal(t, 0.0, queueMetrics.outOfMemory, "Miscount of outOfMemory jobs")
}
