package collector

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestParseCPUs(t *testing.T) {
	assert.Equal(t, CPUs{}, ParseCPUs(""))
	assert.Equal(t, CPUs{alloc: 64, idle: 0, other: 64, total: 128}, ParseCPUs("64/0/64/128"))
	assert.Equal(t, CPUs{alloc: 16, idle: 16, other: 224, total: 256}, ParseCPUs(" 16/16/224/256"))
}

func TestCPUsMetrics(t *testing.T) {
	file, _ := os.Open("fixtures/sinfo/cpus.txt")
	data, _ := io.ReadAll(file)
	cpus := ParseCPUsMetrics(data)

	assert.Equal(t, 5725.0, cpus.alloc, "Miscount of alloc CPUs")
	assert.Equal(t, 877.0, cpus.idle, "Miscount of idle CPUs")
	assert.Equal(t, 34.0, cpus.other, "Miscount of other CPUs")
	assert.Equal(t, 6636.0, cpus.total, "Miscount of total CPUs")
}
