package collector

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseGenericResources(t *testing.T) {
	assert.Equal(t, []GenericResource{}, ParseGenericResources(""))
	assert.Equal(t, []GenericResource{{gresType: "gpu", name: "", count: 8}}, ParseGenericResources("gpu:8"))
	assert.Equal(t, []GenericResource{{gresType: "gpu", name: "tesla", count: 2}, {gresType: "gpu", name: "kepler", count: 2}, {gresType: "mps", count: 400}}, ParseGenericResources("gpu:tesla:2,gpu:kepler:2,mps:400"))
	assert.Equal(t, []GenericResource{{gresType: "gpu", name: "a100m40", count: 8}}, ParseGenericResources("gpu:a100m40:8"))
	assert.Equal(t, []GenericResource{{gresType: "gpu", name: "a100m40", count: 4}, {gresType: "gpu", name: "a100m80", count: 2}}, ParseGenericResources("gpu:a100m40:4(IDX:0-3),gpu:a100m80:2(IDX:4-5)"))
	assert.Equal(t, []GenericResource{{gresType: "gpu", count: 8}}, ParseGenericResources("gpu:8(S:0-1)"))
	assert.Equal(t, []GenericResource{{gresType: "gpu", count: 3}}, ParseGenericResources("gpu:(null):3(IDX:0-7)"))
	assert.Equal(t, []GenericResource{{gresType: "gpu", name: "A30", count: 4}, {gresType: "gpu", name: "Q6K", count: 4}}, ParseGenericResources("gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)"))
}

func TestGPUsMetrics(t *testing.T) {
	testDataPaths, _ := filepath.Glob("fixtures/sinfo/slurm-*")
	for _, testDataPath := range testDataPaths {
		slurmVersion := strings.TrimPrefix(testDataPath, "fixtures/sinfo/slurm-")
		t.Logf("slurm-%s", slurmVersion)

		file, _ := os.Open(testDataPath + "/gpus.txt")
		data, _ := io.ReadAll(file)

		metrics := ParseGPUsMetrics(data)
		assert.Equal(t, 41.0, metrics.idle)
		assert.Equal(t, 7.0, metrics.alloc)
		assert.Equal(t, 48.0, metrics.total)
	}
}
