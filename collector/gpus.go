/*
	Copyright 2020 Joeri Hermans, Victor Penso, Matteo Dessalvi
	Copyright 2022 Iztok Lebar Bajec
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
	"github.com/go-kit/log"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type GPUsMetrics struct {
	alloc       float64
	idle        float64
	other       float64
	total       float64
	utilization float64
}

func ParseAllocatedGPUs(data []byte) float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format="Nodes: ,GresUsed:" --state=allocated
	// 3 gpu:2                                       # slurm>=20.11.8
	// 1 gpu:(null):3(IDX:0-7)                       # slurm 21.08.5
	// 13 gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)      # slurm 21.08.5

	sinfo_lines := string(data)
	re := regexp.MustCompile(`gpu:(\(null\)|[^:(]*):?([0-9]+)(\([^)]*\))?`)
	if len(sinfo_lines) > 0 {
		for _, line := range strings.Split(sinfo_lines, "\n") {
			// log.info(line)
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_active_gpus := strings.Fields(line)[1]
				num_node_active_gpus := 0.0
				for _, node_active_gpus_type := range strings.Split(node_active_gpus, ",") {
					if strings.Contains(node_active_gpus_type, "gpu:") {
						node_active_gpus_type = re.FindStringSubmatch(node_active_gpus_type)[2]
						num_node_active_gpus_type, _ := strconv.ParseFloat(node_active_gpus_type, 64)
						num_node_active_gpus += num_node_active_gpus_type
					}
				}
				num_gpus += num_nodes * num_node_active_gpus
			}
		}
	}

	return num_gpus
}

func ParseIdleGPUs(data []byte) float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format="Nodes: ,Gres: ,GresUsed:" --state=idle,allocated
	// 3 gpu:4 gpu:2                                       																# slurm 20.11.8
	// 1 gpu:8(S:0-1) gpu:(null):3(IDX:0-7)                       												# slurm 21.08.5
	// 13 gpu:A30:4(S:0-1),gpu:Q6K:40(S:0-1) gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)       	# slurm 21.08.5

	sinfo_lines := string(data)
	re := regexp.MustCompile(`gpu:(\(null\)|[^:(]*):?([0-9]+)(\([^)]*\))?`)
	if len(sinfo_lines) > 0 {
		for _, line := range strings.Split(sinfo_lines, "\n") {
			// log.info(line)
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_gpus := strings.Fields(line)[1]
				num_node_gpus := 0.0
				for _, node_gpus_type := range strings.Split(node_gpus, ",") {
					if strings.Contains(node_gpus_type, "gpu:") {
						node_gpus_type = re.FindStringSubmatch(node_gpus_type)[2]
						num_node_gpus_type, _ := strconv.ParseFloat(node_gpus_type, 64)
						num_node_gpus += num_node_gpus_type
					}
				}
				num_node_active_gpus := 0.0
				node_active_gpus := strings.Fields(line)[2]
				for _, node_active_gpus_type := range strings.Split(node_active_gpus, ",") {
					if strings.Contains(node_active_gpus_type, "gpu:") {
						node_active_gpus_type = re.FindStringSubmatch(node_active_gpus_type)[2]
						num_node_active_gpus_type, _ := strconv.ParseFloat(node_active_gpus_type, 64)
						num_node_active_gpus += num_node_active_gpus_type
					}
				}
				num_gpus += num_nodes * (num_node_gpus - num_node_active_gpus)
			}
		}
	}

	return num_gpus
}

func ParseTotalGPUs(data []byte) float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format="Nodes: ,Gres:"
	// 3 gpu:4                                       	# slurm 20.11.8
	// 1 gpu:8(S:0-1)                                	# slurm 21.08.5
	// 13 gpu:A30:4(S:0-1),gpu:Q6K:40(S:0-1)        	# slurm 21.08.5

	sinfo_lines := string(data)
	re := regexp.MustCompile(`gpu:(\(null\)|[^:(]*):?([0-9]+)(\([^)]*\))?`)
	if len(sinfo_lines) > 0 {
		for _, line := range strings.Split(sinfo_lines, "\n") {
			// log.Info(line)
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_gpus := strings.Fields(line)[1]
				num_node_gpus := 0.0
				for _, node_gpus_type := range strings.Split(node_gpus, ",") {
					if strings.Contains(node_gpus_type, "gpu:") {
						node_gpus_type = re.FindStringSubmatch(node_gpus_type)[2]
						num_node_gpus_type, _ := strconv.ParseFloat(node_gpus_type, 64)
						num_node_gpus += num_node_gpus_type
					}
				}
				num_gpus += num_nodes * num_node_gpus
			}
		}
	}

	return num_gpus
}

func GetGPUsMetrics() *GPUsMetrics {
	var gm GPUsMetrics
	total_gpus, _ := GetTotalGPUs()
	allocated_gpus, _ := GetAllocatedGPUs()
	idle_gpus, _ := GetIdleGPUs()
	other_gpus := total_gpus - allocated_gpus - idle_gpus
	gm.alloc = allocated_gpus
	gm.idle = idle_gpus
	gm.other = other_gpus
	gm.total = total_gpus
	gm.utilization = allocated_gpus / total_gpus
	return &gm
}

func GetAllocatedGPUs() (float64, error) {
	out, err := RunCommand("sinfo", "-a", "-h", "--Format=Nodes: ,GresUsed:", "--state=allocated")
	if err != nil {
		return 0, err
	}

	return ParseAllocatedGPUs(out), nil
}

func GetIdleGPUs() (float64, error) {
	out, err := RunCommand("sinfo", "-a", "-h", "--Format=Nodes: ,Gres: ,GresUsed:", "--state=idle,allocated")
	if err != nil {
		return 0, err
	}

	return ParseIdleGPUs(out), nil
}

func GetTotalGPUs() (float64, error) {
	out, err := RunCommand("sinfo", "-a", "-h", "--Format=Nodes: ,Gres:")
	if err != nil {
		return 0, err
	}

	return ParseTotalGPUs(out), nil
}

type GPUsCollector struct {
	alloc       *prometheus.Desc
	idle        *prometheus.Desc
	other       *prometheus.Desc
	total       *prometheus.Desc
	utilization *prometheus.Desc
	logger      log.Logger
}

func init() {
	registerCollector("gpus", defaultEnabled, NewGPUsCollector)
}

func NewGPUsCollector(logger log.Logger) (Collector, error) {
	return &GPUsCollector{
		logger:      logger,
		alloc:       prometheus.NewDesc("slurm_gpus_alloc", "Allocated GPUs", nil, nil),
		idle:        prometheus.NewDesc("slurm_gpus_idle", "Idle GPUs", nil, nil),
		other:       prometheus.NewDesc("slurm_gpus_other", "Other GPUs", nil, nil),
		total:       prometheus.NewDesc("slurm_gpus_total", "Total GPUs", nil, nil),
		utilization: prometheus.NewDesc("slurm_gpus_utilization", "Total GPU utilization", nil, nil),
	}, nil
}

func (cc *GPUsCollector) Collect(ch chan<- prometheus.Metric) error {
	cm := GetGPUsMetrics()
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
	ch <- prometheus.MustNewConstMetric(cc.utilization, prometheus.GaugeValue, cm.utilization)

	return nil
}
