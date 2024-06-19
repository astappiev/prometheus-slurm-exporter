/*
	Copyright 2020 Victor Penso
	Copyright 2024 Oleh Astappiev

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package collector

import (
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type PartitionMetrics struct {
	cpu         CPUs
	jobsPending float64
	jobsRunning float64
}

func ParsePartitionMetrics(input []byte, runningOutput []byte, pendingOutput []byte) map[string]*PartitionMetrics {
	partitions := make(map[string]*PartitionMetrics)

	for _, line := range SplitLines(input) {
		if strings.Contains(line, ",") {
			parts := strings.Split(line, ",")

			// name of a partition
			partition := parts[0]
			_, key := partitions[partition]
			if !key {
				partitions[partition] = &PartitionMetrics{}
			}
			partitions[partition].cpu = ParseCPUs(parts[1])
		}
	}

	// get list of pending jobs by partition name
	for _, partition := range SplitLines(pendingOutput) {
		// accumulate the number of pending jobs
		_, key := partitions[partition]
		if key {
			partitions[partition].jobsPending += 1
		}
	}

	// get list of running jobs by partition name
	for _, partition := range SplitLines(runningOutput) {
		// accumulate the number of running jobs
		_, key := partitions[partition]
		if key {
			partitions[partition].jobsRunning += 1
		}
	}

	return partitions
}

type PartitionCollector struct {
	allocated *prometheus.Desc
	idle      *prometheus.Desc
	other     *prometheus.Desc
	pending   *prometheus.Desc
	running   *prometheus.Desc
	total     *prometheus.Desc
	logger    log.Logger
}

func init() {
	registerCollector("partition", defaultEnabled, NewPartitionCollector)
}

func NewPartitionCollector(logger log.Logger) (Collector, error) {
	return &PartitionCollector{
		logger:    logger,
		allocated: prometheus.NewDesc("slurm_partition_cpus_allocated", "Allocated CPUs for partition", []string{"partition"}, nil),
		idle:      prometheus.NewDesc("slurm_partition_cpus_idle", "Idle CPUs for partition", []string{"partition"}, nil),
		other:     prometheus.NewDesc("slurm_partition_cpus_other", "Other CPUs for partition", []string{"partition"}, nil),
		pending:   prometheus.NewDesc("slurm_partition_jobs_pending", "Pending jobs for partition", []string{"partition"}, nil),
		running:   prometheus.NewDesc("slurm_partition_jobs_running", "Running jobs for partition", []string{"partition"}, nil),
		total:     prometheus.NewDesc("slurm_partition_cpus_total", "Total CPUs for partition", []string{"partition"}, nil),
	}, nil
}

func (pc *PartitionCollector) Collect(ch chan<- prometheus.Metric) error {
	sinfoOutput, err := RunCommand("sinfo", "-h", "-o%R,%C")
	if err != nil {
		return err
	}
	squeueRunningOutput, err := RunCommand("squeue", "-a", "-r", "-h", "-o%P", "--states=RUNNING")
	if err != nil {
		return err
	}
	squeuePendingOutput, err := RunCommand("squeue", "-a", "-r", "-h", "-o%P", "--states=PENDING")
	if err != nil {
		return err
	}

	pm := ParsePartitionMetrics(sinfoOutput, squeueRunningOutput, squeuePendingOutput)
	for p := range pm {
		if pm[p].cpu.alloc > 0 {
			ch <- prometheus.MustNewConstMetric(pc.allocated, prometheus.GaugeValue, pm[p].cpu.alloc, p)
		}
		if pm[p].cpu.idle > 0 {
			ch <- prometheus.MustNewConstMetric(pc.idle, prometheus.GaugeValue, pm[p].cpu.idle, p)
		}
		if pm[p].cpu.other > 0 {
			ch <- prometheus.MustNewConstMetric(pc.other, prometheus.GaugeValue, pm[p].cpu.other, p)
		}
		if pm[p].cpu.total > 0 {
			ch <- prometheus.MustNewConstMetric(pc.total, prometheus.GaugeValue, pm[p].cpu.total, p)
		}
		if pm[p].jobsPending > 0 {
			ch <- prometheus.MustNewConstMetric(pc.pending, prometheus.GaugeValue, pm[p].jobsPending, p)
		}
		if pm[p].jobsRunning > 0 {
			ch <- prometheus.MustNewConstMetric(pc.running, prometheus.GaugeValue, pm[p].jobsRunning, p)
		}
	}

	return nil
}
