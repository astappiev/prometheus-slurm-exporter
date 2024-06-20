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
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type CPUs struct {
	alloc float64
	idle  float64
	other float64
	total float64
}

func ParseCPUs(input string) CPUs {
	var cpus CPUs
	if len(input) > 0 && strings.Contains(input, "/") {
		parts := strings.Split(strings.TrimSpace(input), "/")
		cpus.alloc, _ = strconv.ParseFloat(parts[0], 64)
		cpus.idle, _ = strconv.ParseFloat(parts[1], 64)
		cpus.other, _ = strconv.ParseFloat(parts[2], 64)
		cpus.total, _ = strconv.ParseFloat(parts[3], 64)
	}
	return cpus
}

func ParseCPUsMetrics(input []byte) *CPUs {
	cpu := ParseCPUs(string(input))
	return &cpu
}

type CPUsCollector struct {
	alloc  *prometheus.Desc
	idle   *prometheus.Desc
	other  *prometheus.Desc
	total  *prometheus.Desc
	logger log.Logger
}

func init() {
	registerCollector("cpus", defaultEnabled, NewCPUsCollector)
}

func NewCPUsCollector(logger log.Logger) (Collector, error) {
	return &CPUsCollector{
		logger: logger,
		alloc:  prometheus.NewDesc("slurm_cpus_alloc", "Allocated CPUs", nil, nil),
		idle:   prometheus.NewDesc("slurm_cpus_idle", "Idle CPUs", nil, nil),
		other:  prometheus.NewDesc("slurm_cpus_other", "Mix CPUs", nil, nil),
		total:  prometheus.NewDesc("slurm_cpus_total", "Total CPUs", nil, nil),
	}, nil
}

func (cc *CPUsCollector) Collect(ch chan<- prometheus.Metric) error {
	out, err := RunCommand("sinfo", "-h", "-a", "-o %C")
	if err != nil {
		return err
	}

	cm := ParseCPUsMetrics(out)
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)

	return nil
}
