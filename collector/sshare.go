/*
	Copyright 2021 Victor Penso
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

type FairShareMetrics struct {
	fairshare float64
}

func ParseFairShareMetrics(input []byte) map[string]*FairShareMetrics {
	accounts := make(map[string]*FairShareMetrics)

	for _, line := range SplitLines(input) {
		if !strings.HasPrefix(line, "  ") {
			if strings.Contains(line, "|") {
				parts := strings.Split(line, "|")
				account := strings.Trim(parts[0], " ")
				_, key := accounts[account]
				if !key {
					accounts[account] = &FairShareMetrics{}
				}
				fairshare, _ := strconv.ParseFloat(parts[1], 64)
				accounts[account].fairshare = fairshare
			}
		}
	}
	return accounts
}

type FairShareCollector struct {
	fairshare *prometheus.Desc
	logger    log.Logger
}

func init() {
	registerCollector("fairshare", defaultDisabled, NewFairShareCollector)
}

func NewFairShareCollector(logger log.Logger) (Collector, error) {
	return &FairShareCollector{
		logger:    logger,
		fairshare: prometheus.NewDesc("slurm_account_fairshare", "FairShare for account", []string{"account"}, nil),
	}, nil
}

func (fsc *FairShareCollector) Collect(ch chan<- prometheus.Metric) error {
	out, err := RunCommand("sshare", "-n", "-P", "-o", "account,fairshare")
	if err != nil {
		return err
	}

	fsm := ParseFairShareMetrics(out)
	for f := range fsm {
		ch <- prometheus.MustNewConstMetric(fsc.fairshare, prometheus.GaugeValue, fsm[f].fairshare, f)
	}

	return nil
}
