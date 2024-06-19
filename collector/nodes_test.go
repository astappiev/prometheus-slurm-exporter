package collector

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestNodesMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("fixtures/sinfo/nodes.txt")
	data, _ := io.ReadAll(file)
	nodes := ParseNodesMetrics(data)

	assert.Equal(t, 1.0, nodes.alloc)
	assert.Equal(t, 0.0, nodes.comp)
	assert.Equal(t, 0.0, nodes.down)
	assert.Equal(t, 7.0, nodes.drain)
	assert.Equal(t, 0.0, nodes.err)
	assert.Equal(t, 0.0, nodes.fail)
	assert.Equal(t, 9.0, nodes.idle)
	assert.Equal(t, 0.0, nodes.maint)
	assert.Equal(t, 20.0, nodes.mix)
	assert.Equal(t, 0.0, nodes.resv)
	assert.Equal(t, 0.0, nodes.plnd)
}
