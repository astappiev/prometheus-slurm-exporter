package collector

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestParseFairShareMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("fixtures/sshare/sshare.txt")
	data, _ := io.ReadAll(file)
	metrics := ParseFairShareMetrics(data)

	assert.Equal(t, 0.0, metrics["ampere"].fairshare)
	assert.Equal(t, 0.0, metrics["volta"].fairshare)
}
