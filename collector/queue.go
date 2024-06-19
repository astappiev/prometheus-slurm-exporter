/*
	Copyright 2017 Victor Penso, Matteo Dessalvi
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

type QueueMetrics struct {
	pending     float64
	pendingDep  float64
	running     float64
	suspended   float64
	cancelled   float64
	completing  float64
	completed   float64
	configuring float64
	failed      float64
	timeout     float64
	preempted   float64
	nodeFail    float64
	outOfMemory float64
}

func ParseQueueMetrics(input []byte) *QueueMetrics {
	var qm QueueMetrics

	for _, line := range SplitLines(input) {
		if strings.Contains(line, ",") {
			parts := strings.Split(line, ",")

			state := parts[1]
			switch state {
			case "PENDING":
				qm.pending++
				if len(parts) > 2 && parts[2] == "Dependency" {
					qm.pendingDep++
				}
			case "RUNNING":
				qm.running++
			case "SUSPENDED":
				qm.suspended++
			case "CANCELLED":
				qm.cancelled++
			case "COMPLETING":
				qm.completing++
			case "COMPLETED":
				qm.completed++
			case "CONFIGURING":
				qm.configuring++
			case "FAILED":
				qm.failed++
			case "TIMEOUT":
				qm.timeout++
			case "PREEMPTED":
				qm.preempted++
			case "NODE_FAIL":
				qm.nodeFail++
			case "OUT_OF_MEMORY":
				qm.outOfMemory++
			}
		}
	}
	return &qm
}

type QueueCollector struct {
	pending     *prometheus.Desc
	pendingDep  *prometheus.Desc
	running     *prometheus.Desc
	suspended   *prometheus.Desc
	cancelled   *prometheus.Desc
	completing  *prometheus.Desc
	completed   *prometheus.Desc
	configuring *prometheus.Desc
	failed      *prometheus.Desc
	timeout     *prometheus.Desc
	preempted   *prometheus.Desc
	nodeFail    *prometheus.Desc
	outOfMemory *prometheus.Desc
	logger      log.Logger
}

func init() {
	registerCollector("queue", defaultEnabled, NewQueueCollector)
}

func NewQueueCollector(logger log.Logger) (Collector, error) {
	return &QueueCollector{
		logger:      logger,
		pending:     prometheus.NewDesc("slurm_queue_pending", "Pending jobs in queue", nil, nil),
		pendingDep:  prometheus.NewDesc("slurm_queue_pending_dependency", "Pending jobs because of dependency in queue", nil, nil),
		running:     prometheus.NewDesc("slurm_queue_running", "Running jobs in the cluster", nil, nil),
		suspended:   prometheus.NewDesc("slurm_queue_suspended", "Suspended jobs in the cluster", nil, nil),
		cancelled:   prometheus.NewDesc("slurm_queue_cancelled", "Cancelled jobs in the cluster", nil, nil),
		completing:  prometheus.NewDesc("slurm_queue_completing", "Completing jobs in the cluster", nil, nil),
		completed:   prometheus.NewDesc("slurm_queue_completed", "Completed jobs in the cluster", nil, nil),
		configuring: prometheus.NewDesc("slurm_queue_configuring", "Configuring jobs in the cluster", nil, nil),
		failed:      prometheus.NewDesc("slurm_queue_failed", "Number of failed jobs", nil, nil),
		timeout:     prometheus.NewDesc("slurm_queue_timeout", "Jobs stopped by timeout", nil, nil),
		preempted:   prometheus.NewDesc("slurm_queue_preempted", "Number of preempted jobs", nil, nil),
		nodeFail:    prometheus.NewDesc("slurm_queue_node_fail", "Number of jobs stopped due to node fail", nil, nil),
		outOfMemory: prometheus.NewDesc("slurm_queue_out_of_memory", "Number of jobs stopped by oomkiller", nil, nil),
	}, nil
}

func (qc *QueueCollector) Collect(ch chan<- prometheus.Metric) error {
	out, err := RunCommand("squeue", "-a", "-r", "-h", "-o %A,%T,%r", "--states=all")
	if err != nil {
		return err
	}

	qm := ParseQueueMetrics(out)
	ch <- prometheus.MustNewConstMetric(qc.pending, prometheus.GaugeValue, qm.pending)
	ch <- prometheus.MustNewConstMetric(qc.pendingDep, prometheus.GaugeValue, qm.pendingDep)
	ch <- prometheus.MustNewConstMetric(qc.running, prometheus.GaugeValue, qm.running)
	ch <- prometheus.MustNewConstMetric(qc.suspended, prometheus.GaugeValue, qm.suspended)
	ch <- prometheus.MustNewConstMetric(qc.cancelled, prometheus.GaugeValue, qm.cancelled)
	ch <- prometheus.MustNewConstMetric(qc.completing, prometheus.GaugeValue, qm.completing)
	ch <- prometheus.MustNewConstMetric(qc.completed, prometheus.GaugeValue, qm.completed)
	ch <- prometheus.MustNewConstMetric(qc.configuring, prometheus.GaugeValue, qm.configuring)
	ch <- prometheus.MustNewConstMetric(qc.failed, prometheus.GaugeValue, qm.failed)
	ch <- prometheus.MustNewConstMetric(qc.timeout, prometheus.GaugeValue, qm.timeout)
	ch <- prometheus.MustNewConstMetric(qc.preempted, prometheus.GaugeValue, qm.preempted)
	ch <- prometheus.MustNewConstMetric(qc.nodeFail, prometheus.GaugeValue, qm.nodeFail)
	ch <- prometheus.MustNewConstMetric(qc.outOfMemory, prometheus.GaugeValue, qm.outOfMemory)

	return nil
}
