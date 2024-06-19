package collector

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGPUsMetrics(t *testing.T) {
	test_data_paths, _ := filepath.Glob("fixtures/sinfo/slurm-*")
	for _, test_data_path := range test_data_paths {
		slurm_version := strings.TrimPrefix(test_data_path, "fixtures/sinfo/slurm-")
		t.Logf("slurm-%s", slurm_version)

		// Read the input data from a file
		file, err := os.Open(test_data_path + "/sinfo_gpus_allocated.txt")
		if err != nil {
			t.Fatalf("Can not open test data: %v", err)
		}
		data, _ := io.ReadAll(file)
		metrics := ParseAllocatedGPUs(data)
		t.Logf("Allocated: %+v", metrics)

		// Read the input data from a file
		file, err = os.Open(test_data_path + "/sinfo_gpus_idle.txt")
		if err != nil {
			t.Fatalf("Can not open test data: %v", err)
		}
		data, _ = io.ReadAll(file)
		metrics = ParseIdleGPUs(data)
		t.Logf("Idle: %+v", metrics)

		// Read the input data from a file
		file, err = os.Open(test_data_path + "/sinfo_gpus_total.txt")
		if err != nil {
			t.Fatalf("Can not open test data: %v", err)
		}
		data, _ = io.ReadAll(file)
		metrics = ParseTotalGPUs(data)
		t.Logf("Total: %+v", metrics)
	}
}
