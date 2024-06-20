package collector

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("fixtures/sinfo/node.txt")
	data, _ := io.ReadAll(file)
	metrics := ParseNodeMetrics(data)

	assert.Contains(t, metrics, "gpunode05")
	assert.Equal(t, float64(0), metrics["gpunode05"].memAlloc)
	assert.Equal(t, float64(515500), metrics["gpunode05"].memTotal)
	assert.Equal(t, float64(60), metrics["gpunode05"].cpu.alloc)
	assert.Equal(t, float64(68), metrics["gpunode05"].cpu.idle)
	assert.Equal(t, float64(0), metrics["gpunode05"].cpu.other)
	assert.Equal(t, float64(128), metrics["gpunode05"].cpu.total)
	assert.Equal(t, "a100m40", metrics["gpunode05"].gres[0].name)
	assert.Equal(t, float64(4), metrics["gpunode05"].gres[0].count)
	assert.Equal(t, "a100m80", metrics["gpunode05"].gres[1].name)
	assert.Equal(t, float64(4), metrics["gpunode05"].gres[1].count)
	assert.Equal(t, "a100m40", metrics["gpunode05"].gresUsed[0].name)
	assert.Equal(t, float64(4), metrics["gpunode05"].gresUsed[0].count)
	assert.Equal(t, "a100m80", metrics["gpunode05"].gresUsed[1].name)
	assert.Equal(t, float64(4), metrics["gpunode05"].gresUsed[1].count)
}
