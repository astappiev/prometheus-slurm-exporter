/*
	Copyright 2021 Chris Read
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
	"sort"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

// NodeMetrics stores metrics for each node
type NodeMetrics struct {
	memAlloc   float64
	memTotal   float64
	cpu        CPUs
	tres       []TRES
	tresUsed   []TRES
	nodeStatus string
}

type TRES struct {
	kind  string
	name  string
	count float64
}

func ParseTres(input string) []TRES {
	results := make([]TRES, 0)

	// remove parenthesis and all content inside
	for strings.Contains(input, "(") {
		openIndex := strings.Index(input, "(")
		closeIndex := strings.Index(input, ")")
		input = input[:openIndex] + input[closeIndex+1:]
	}

	for _, single := range strings.Split(input, ",") {
		tresInfo := strings.Split(single, ":")
		tresCount, _ := strconv.ParseFloat(tresInfo[2], 64)
		results = append(results, TRES{kind: tresInfo[0], name: tresInfo[1], count: tresCount})
	}

	return results
}

// ParseNodeMetrics takes the output of sinfo with node data
func ParseNodeMetrics(input []byte) map[string]*NodeMetrics {
	nodes := make(map[string]*NodeMetrics)

	lines := SplitLines(input)
	// Sort and remove all the duplicates from the 'sinfo' output
	sort.Strings(lines)
	linesUniq := RemoveDuplicates(lines)

	for _, line := range linesUniq {
		node := strings.Fields(line)
		nodeName := node[0]
		nodeStatus := node[4] // mixed, allocated, etc.

		nodes[nodeName] = &NodeMetrics{}

		nodes[nodeName].cpu = ParseCPUs(node[3])
		memAlloc, _ := strconv.ParseFloat(node[1], 64)
		memTotal, _ := strconv.ParseFloat(node[2], 64)

		if len(node) >= 6 && node[5] != "(null)" && len(node[5]) > 0 {
			nodes[nodeName].tres = ParseTres(node[5])
			nodes[nodeName].tresUsed = ParseTres(node[6])
		}

		nodes[nodeName].memAlloc = memAlloc
		nodes[nodeName].memTotal = memTotal
		nodes[nodeName].nodeStatus = nodeStatus
	}

	return nodes
}

type NodeCollector struct {
	cpuAlloc *prometheus.Desc
	cpuIdle  *prometheus.Desc
	cpuOther *prometheus.Desc
	cpuTotal *prometheus.Desc
	memAlloc *prometheus.Desc
	memTotal *prometheus.Desc
	gpuAlloc *prometheus.Desc
	gpuTotal *prometheus.Desc
	logger   log.Logger
}

func init() {
	registerCollector("node", defaultEnabled, NewNodeCollector)
}

// NewNodeCollector creates a Prometheus collector to keep all our stats in
// It returns a set of collections for consumption
func NewNodeCollector(logger log.Logger) (Collector, error) {
	return &NodeCollector{
		logger:   logger,
		cpuAlloc: prometheus.NewDesc("slurm_node_cpu_alloc", "Allocated CPUs per node", []string{"node", "status"}, nil),
		cpuIdle:  prometheus.NewDesc("slurm_node_cpu_idle", "Idle CPUs per node", []string{"node", "status"}, nil),
		cpuOther: prometheus.NewDesc("slurm_node_cpu_other", "Other CPUs per node", []string{"node", "status"}, nil),
		cpuTotal: prometheus.NewDesc("slurm_node_cpu_total", "Total CPUs per node", []string{"node", "status"}, nil),
		memAlloc: prometheus.NewDesc("slurm_node_mem_alloc", "Allocated memory per node", []string{"node", "status"}, nil),
		memTotal: prometheus.NewDesc("slurm_node_mem_total", "Total memory per node", []string{"node", "status"}, nil),
		gpuAlloc: prometheus.NewDesc("slurm_node_gpu_alloc", "Allocated GPUs per node", []string{"node", "status", "gputype"}, nil),
		gpuTotal: prometheus.NewDesc("slurm_node_gpu_total", "Total GPUs per node", []string{"node", "status", "gputype"}, nil),
	}, nil
}

func (c *NodeCollector) Collect(ch chan<- prometheus.Metric) error {
	out, err := RunCommand("sinfo", "-h", "-a", "-N", "-O", "NodeList: ,AllocMem: ,Memory: ,CPUsState: ,StateLong: ,Gres: ,Gresused:")
	if err != nil {
		return err
	}

	nodes := ParseNodeMetrics(out)
	for node := range nodes {
		ch <- prometheus.MustNewConstMetric(c.cpuAlloc, prometheus.GaugeValue, nodes[node].cpu.alloc, node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(c.cpuIdle, prometheus.GaugeValue, nodes[node].cpu.idle, node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(c.cpuOther, prometheus.GaugeValue, nodes[node].cpu.other, node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(c.cpuTotal, prometheus.GaugeValue, nodes[node].cpu.total, node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(c.memAlloc, prometheus.GaugeValue, nodes[node].memAlloc, node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(c.memTotal, prometheus.GaugeValue, nodes[node].memTotal, node, nodes[node].nodeStatus)
		if nodes[node].tresUsed != nil && len(nodes[node].tresUsed) != 0 {
			for _, tres := range nodes[node].tresUsed {
				ch <- prometheus.MustNewConstMetric(c.gpuAlloc, prometheus.GaugeValue, tres.count, node, nodes[node].nodeStatus, tres.name)
			}
		}
		if nodes[node].tres != nil && len(nodes[node].tres) != 0 {
			for _, tres := range nodes[node].tres {
				ch <- prometheus.MustNewConstMetric(c.gpuTotal, prometheus.GaugeValue, tres.count, node, nodes[node].nodeStatus, tres.name)
			}
		}
	}

	return nil
}
