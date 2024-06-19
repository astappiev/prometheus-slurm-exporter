package collector

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestParseUsersMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("fixtures/squeue/user.txt")
	data, _ := io.ReadAll(file)
	users := ParseUserMetrics(data)

	assert.Equal(t, 31.0, users["user2"].jobsPending, "Miscount of pending user jobs")
	assert.Equal(t, 8.0, users["user2"].jobsRunning, "Miscount of running user jobs")
	assert.Equal(t, 0.0, users["user2"].jobsSuspended, "Miscount of suspended user jobs")
	assert.Equal(t, 32.0, users["user2"].cpusRunning, "Miscount of running user CPUs")
	assert.Equal(t, 124.0, users["user2"].cpusPending, "Miscount of pending user CPUs")
	assert.Equal(t, 2.74877906944e+11, users["user2"].memRunning, "Miscount of running user Memory")
}
