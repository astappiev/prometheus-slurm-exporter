package collector

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestSchedulerMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("fixtures/sdiag/sdiag.txt")
	data, _ := io.ReadAll(file)
	schedulerMetrics := ParseSchedulerMetrics(data)

	assert.Equal(t, 2.0, schedulerMetrics.threads)
	assert.Equal(t, 0.0, schedulerMetrics.queueSize)
	assert.Equal(t, 0.0, schedulerMetrics.dbdQueueSize)
	assert.Equal(t, 2291.0, schedulerMetrics.lastCycle)
	assert.Equal(t, 2498.0, schedulerMetrics.meanCycle)
	assert.Equal(t, 1.0, schedulerMetrics.cyclePerMinute)
	assert.Equal(t, 5909.0, schedulerMetrics.backfillLastCycle)
	assert.Equal(t, 4799.0, schedulerMetrics.backfillMeanCycle)
	assert.Equal(t, 37.0, schedulerMetrics.backfillDepthMean)
	assert.Equal(t, 155.0, schedulerMetrics.totalBackfilledJobsSinceStart)
	assert.Equal(t, 6.0, schedulerMetrics.totalBackfilledJobsSinceCycle)
	assert.Equal(t, 0.0, schedulerMetrics.totalBackfilledHeterogeneous)
}
