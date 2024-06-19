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
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type UserJobMetrics struct {
	jobsPending   float64
	cpusPending   float64
	jobsRunning   float64
	cpusRunning   float64
	memRunning    float64
	jobsSuspended float64
}

func ParseMemory(input string) float64 {
	if len(input) == 0 || input == "0" {
		return 0
	}

	reg := `^(\d+)([KMGT])$`
	r := regexp.MustCompile(reg)
	matchs := r.FindStringSubmatch(string(input))
	num, _ := strconv.Atoi(matchs[1])
	unit := matchs[2]
	ret := 0
	if "K" == unit {
		ret = num * 1024
	} else if "M" == unit {
		ret = num * 1024 * 1024
	} else if "G" == unit {
		ret = num * 1024 * 1024 * 1024
	} else if "T" == unit {
		ret = num * 1024 * 1024 * 1024 * 1024
	} else {
		ret = num
	}
	return float64(ret)
}

func ParseUserMetrics(input []byte) map[string]*UserJobMetrics {
	users := make(map[string]*UserJobMetrics)

	var (
		pending   = regexp.MustCompile(`^pending`)
		running   = regexp.MustCompile(`^running`)
		suspended = regexp.MustCompile(`^suspended`)
	)

	for _, line := range SplitLines(input) {
		if strings.Contains(line, "|") {
			parts := strings.Split(line, "|")
			user := parts[1]
			_, key := users[user]
			if !key {
				users[user] = &UserJobMetrics{}
			}
			state := parts[2]
			state = strings.ToLower(state)
			cpus, _ := strconv.ParseFloat(parts[3], 64)
			mem := ParseMemory(parts[4])
			switch {
			case pending.MatchString(state) == true:
				users[user].jobsPending++
				users[user].cpusPending += cpus
			case running.MatchString(state) == true:
				users[user].jobsRunning++
				users[user].cpusRunning += cpus
				users[user].memRunning += mem
			case suspended.MatchString(state) == true:
				users[user].jobsSuspended++
			}
		}
	}
	return users
}

type UserCollector struct {
	jobsPending   *prometheus.Desc
	cpusPending   *prometheus.Desc
	jobsRunning   *prometheus.Desc
	cpusRunning   *prometheus.Desc
	memRunning    *prometheus.Desc
	jobsSuspended *prometheus.Desc
	logger        log.Logger
}

func init() {
	registerCollector("user", defaultEnabled, NewUserCollector)
}

func NewUserCollector(logger log.Logger) (Collector, error) {
	return &UserCollector{
		logger:        logger,
		jobsPending:   prometheus.NewDesc("slurm_user_jobs_pending", "Pending jobs for user", []string{"user"}, nil),
		cpusPending:   prometheus.NewDesc("slurm_user_cpus_pending", "Pending jobs for user", []string{"user"}, nil),
		jobsRunning:   prometheus.NewDesc("slurm_user_jobs_running", "Running jobs for user", []string{"user"}, nil),
		cpusRunning:   prometheus.NewDesc("slurm_user_cpus_running", "Running cpus for user", []string{"user"}, nil),
		memRunning:    prometheus.NewDesc("slurm_user_mem_running", "Running mem for user", []string{"user"}, nil),
		jobsSuspended: prometheus.NewDesc("slurm_user_jobs_suspended", "Suspended jobs for user", []string{"user"}, nil),
	}, nil
}

func (uc *UserCollector) Collect(ch chan<- prometheus.Metric) error {
	out, err := RunCommand("squeue", "-a", "-r", "-h", "-o %A|%u|%T|%C|%m")
	if err != nil {
		return err
	}

	um := ParseUserMetrics(out)
	for u := range um {
		if um[u].jobsPending > 0 {
			ch <- prometheus.MustNewConstMetric(uc.jobsPending, prometheus.GaugeValue, um[u].jobsPending, u)
		}
		if um[u].cpusPending > 0 {
			ch <- prometheus.MustNewConstMetric(uc.cpusPending, prometheus.GaugeValue, um[u].cpusPending, u)
		}
		if um[u].jobsRunning > 0 {
			ch <- prometheus.MustNewConstMetric(uc.jobsRunning, prometheus.GaugeValue, um[u].jobsRunning, u)
		}
		if um[u].cpusRunning > 0 {
			ch <- prometheus.MustNewConstMetric(uc.cpusRunning, prometheus.GaugeValue, um[u].cpusRunning, u)
		}
		if um[u].memRunning > 0 {
			ch <- prometheus.MustNewConstMetric(uc.memRunning, prometheus.GaugeValue, um[u].memRunning, u)
		}
		if um[u].jobsSuspended > 0 {
			ch <- prometheus.MustNewConstMetric(uc.jobsSuspended, prometheus.GaugeValue, um[u].jobsSuspended, u)
		}
	}

	return nil
}
