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

type JobMetrics struct {
	pending     float64
	pendingCpus float64
	running     float64
	runningCpus float64
	suspended   float64
}

func ParseAccountMetrics(input []byte) map[string]*JobMetrics {
	accounts := make(map[string]*JobMetrics)

	var (
		pending   = regexp.MustCompile(`^pending`)
		running   = regexp.MustCompile(`^running`)
		suspended = regexp.MustCompile(`^suspended`)
	)

	for _, line := range SplitLines(input) {
		if strings.Contains(line, "|") {
			parts := strings.Split(line, "|")

			account := parts[1]
			_, key := accounts[account]
			if !key {
				accounts[account] = &JobMetrics{}
			}
			state := strings.ToLower(parts[2])
			cpus, _ := strconv.ParseFloat(parts[3], 64)

			switch {
			case pending.MatchString(state) == true:
				accounts[account].pending++
				accounts[account].pendingCpus += cpus
			case running.MatchString(state) == true:
				accounts[account].running++
				accounts[account].runningCpus += cpus
			case suspended.MatchString(state) == true:
				accounts[account].suspended++
			}
		}
	}
	return accounts
}

type AccountCollector struct {
	pending     *prometheus.Desc
	pendingCpus *prometheus.Desc
	running     *prometheus.Desc
	runningCpus *prometheus.Desc
	suspended   *prometheus.Desc
	logger      log.Logger
}

func init() {
	registerCollector("account", defaultEnabled, NewAccountCollector)
}

func NewAccountCollector(logger log.Logger) (Collector, error) {
	return &AccountCollector{
		logger:      logger,
		pending:     prometheus.NewDesc("slurm_account_jobs_pending", "Pending jobs for account", []string{"account"}, nil),
		pendingCpus: prometheus.NewDesc("slurm_account_cpus_pending", "Pending jobs for account", []string{"account"}, nil),
		running:     prometheus.NewDesc("slurm_account_jobs_running", "Running jobs for account", []string{"account"}, nil),
		runningCpus: prometheus.NewDesc("slurm_account_cpus_running", "Running cpus for account", []string{"account"}, nil),
		suspended:   prometheus.NewDesc("slurm_account_jobs_suspended", "Suspended jobs for account", []string{"account"}, nil),
	}, nil
}

func (ac *AccountCollector) Collect(ch chan<- prometheus.Metric) error {
	out, err := RunCommand("squeue", "-a", "-r", "-h", "-o %A|%a|%T|%C")
	if err != nil {
		return err
	}

	am := ParseAccountMetrics(out)
	for a := range am {
		if am[a].pending > 0 {
			ch <- prometheus.MustNewConstMetric(ac.pending, prometheus.GaugeValue, am[a].pending, a)
		}
		if am[a].pendingCpus > 0 {
			ch <- prometheus.MustNewConstMetric(ac.pendingCpus, prometheus.GaugeValue, am[a].pendingCpus, a)
		}
		if am[a].running > 0 {
			ch <- prometheus.MustNewConstMetric(ac.running, prometheus.GaugeValue, am[a].running, a)
		}
		if am[a].runningCpus > 0 {
			ch <- prometheus.MustNewConstMetric(ac.runningCpus, prometheus.GaugeValue, am[a].runningCpus, a)
		}
		if am[a].suspended > 0 {
			ch <- prometheus.MustNewConstMetric(ac.suspended, prometheus.GaugeValue, am[a].suspended, a)
		}
	}

	return nil
}
