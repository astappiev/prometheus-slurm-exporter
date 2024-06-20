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
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type SchedulerMetrics struct {
	threads                       float64
	queueSize                     float64
	dbdQueueSize                  float64
	lastCycle                     float64
	meanCycle                     float64
	cyclePerMinute                float64
	backfillLastCycle             float64
	backfillMeanCycle             float64
	backfillDepthMean             float64
	totalBackfilledJobsSinceStart float64
	totalBackfilledJobsSinceCycle float64
	totalBackfilledHeterogeneous  float64
}

func ParseSchedulerMetrics(input []byte) *SchedulerMetrics {
	var (
		st  = regexp.MustCompile(`^Server thread`)
		qs  = regexp.MustCompile(`^Agent queue`)
		dbd = regexp.MustCompile(`^DBD Agent`)
		lc  = regexp.MustCompile(`^[\s]+Last cycle$`)
		mc  = regexp.MustCompile(`^[\s]+Mean cycle$`)
		cpm = regexp.MustCompile(`^[\s]+Cycles per`)
		dpm = regexp.MustCompile(`^[\s]+Depth Mean$`)
		tbs = regexp.MustCompile(`^[\s]+Total backfilled jobs \(since last slurm start\)`)
		tbc = regexp.MustCompile(`^[\s]+Total backfilled jobs \(since last stats cycle start\)`)
		tbh = regexp.MustCompile(`^[\s]+Total backfilled heterogeneous job components`)
	)

	var sm SchedulerMetrics
	lines := SplitLines(input)
	// Guard variables to check for string repetitions in the output of sdiag
	// (two occurrences of the following strings: 'Last cycle', 'Mean cycle')
	lcCount := 0
	mcCount := 0
	for _, line := range lines {
		if strings.Contains(line, ":") {
			state := strings.Split(line, ":")[0]
			switch {
			case st.MatchString(state) == true:
				sm.threads, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
			case qs.MatchString(state) == true:
				sm.queueSize, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
			case dbd.MatchString(state) == true:
				sm.dbdQueueSize, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
			case lc.MatchString(state) == true:
				if lcCount == 0 {
					sm.lastCycle, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
					lcCount = 1
				} else {
					sm.backfillLastCycle, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
				}
			case mc.MatchString(state) == true:
				if mcCount == 0 {
					sm.meanCycle, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
					mcCount = 1
				} else {
					sm.backfillMeanCycle, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
				}
			case cpm.MatchString(state) == true:
				sm.cyclePerMinute, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
			case dpm.MatchString(state) == true:
				sm.backfillDepthMean, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
			case tbs.MatchString(state) == true:
				sm.totalBackfilledJobsSinceStart, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
			case tbc.MatchString(state) == true:
				sm.totalBackfilledJobsSinceCycle, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
			case tbh.MatchString(state) == true:
				sm.totalBackfilledHeterogeneous, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
			}
		}
	}
	return &sm
}

type SchedulerCollector struct {
	threads                       *prometheus.Desc
	queueSize                     *prometheus.Desc
	dbdQueueSize                  *prometheus.Desc
	lastCycle                     *prometheus.Desc
	meanCycle                     *prometheus.Desc
	cyclePerMinute                *prometheus.Desc
	backfillLastCycle             *prometheus.Desc
	backfillMeanCycle             *prometheus.Desc
	backfillDepthMean             *prometheus.Desc
	totalBackfilledJobsSinceStart *prometheus.Desc
	totalBackfilledJobsSinceCycle *prometheus.Desc
	totalBackfilledHeterogeneous  *prometheus.Desc
	logger                        log.Logger
}

func init() {
	registerCollector("scheduler", defaultEnabled, NewSchedulerCollector)
}

func NewSchedulerCollector(logger log.Logger) (Collector, error) {
	return &SchedulerCollector{
		logger:                        logger,
		threads:                       prometheus.NewDesc("slurm_scheduler_threads", "Information provided by the Slurm sdiag command, number of scheduler threads ", nil, nil),
		queueSize:                     prometheus.NewDesc("slurm_scheduler_queue_size", "Information provided by the Slurm sdiag command, length of the scheduler queue", nil, nil),
		dbdQueueSize:                  prometheus.NewDesc("slurm_scheduler_dbd_queue_size", "Information provided by the Slurm sdiag command, length of the DBD agent queue", nil, nil),
		lastCycle:                     prometheus.NewDesc("slurm_scheduler_last_cycle", "Information provided by the Slurm sdiag command, scheduler last cycle time in (microseconds)", nil, nil),
		meanCycle:                     prometheus.NewDesc("slurm_scheduler_mean_cycle", "Information provided by the Slurm sdiag command, scheduler mean cycle time in (microseconds)", nil, nil),
		cyclePerMinute:                prometheus.NewDesc("slurm_scheduler_cycle_per_minute", "Information provided by the Slurm sdiag command, number scheduler cycles per minute", nil, nil),
		backfillLastCycle:             prometheus.NewDesc("slurm_scheduler_backfill_last_cycle", "Information provided by the Slurm sdiag command, scheduler backfill last cycle time in (microseconds)", nil, nil),
		backfillMeanCycle:             prometheus.NewDesc("slurm_scheduler_backfill_mean_cycle", "Information provided by the Slurm sdiag command, scheduler backfill mean cycle time in (microseconds)", nil, nil),
		backfillDepthMean:             prometheus.NewDesc("slurm_scheduler_backfill_depth_mean", "Information provided by the Slurm sdiag command, scheduler backfill mean depth", nil, nil),
		totalBackfilledJobsSinceStart: prometheus.NewDesc("slurm_scheduler_backfilled_jobs_since_start_total", "Information provided by the Slurm sdiag command, number of jobs started thanks to backfilling since last slurm start", nil, nil),
		totalBackfilledJobsSinceCycle: prometheus.NewDesc("slurm_scheduler_backfilled_jobs_since_cycle_total", "Information provided by the Slurm sdiag command, number of jobs started thanks to backfilling since last time stats where reset", nil, nil),
		totalBackfilledHeterogeneous:  prometheus.NewDesc("slurm_scheduler_backfilled_heterogeneous_total", "Information provided by the Slurm sdiag command, number of heterogeneous job components started thanks to backfilling since last Slurm start", nil, nil),
	}, nil
}

func (sc *SchedulerCollector) Collect(ch chan<- prometheus.Metric) error {
	out, err := RunCommand("sdiag")
	if err != nil {
		return err
	}

	sm := ParseSchedulerMetrics(out)
	ch <- prometheus.MustNewConstMetric(sc.threads, prometheus.GaugeValue, sm.threads)
	ch <- prometheus.MustNewConstMetric(sc.queueSize, prometheus.GaugeValue, sm.queueSize)
	ch <- prometheus.MustNewConstMetric(sc.dbdQueueSize, prometheus.GaugeValue, sm.dbdQueueSize)
	ch <- prometheus.MustNewConstMetric(sc.lastCycle, prometheus.GaugeValue, sm.lastCycle)
	ch <- prometheus.MustNewConstMetric(sc.meanCycle, prometheus.GaugeValue, sm.meanCycle)
	ch <- prometheus.MustNewConstMetric(sc.cyclePerMinute, prometheus.GaugeValue, sm.cyclePerMinute)
	ch <- prometheus.MustNewConstMetric(sc.backfillLastCycle, prometheus.GaugeValue, sm.backfillLastCycle)
	ch <- prometheus.MustNewConstMetric(sc.backfillMeanCycle, prometheus.GaugeValue, sm.backfillMeanCycle)
	ch <- prometheus.MustNewConstMetric(sc.backfillDepthMean, prometheus.GaugeValue, sm.backfillDepthMean)
	ch <- prometheus.MustNewConstMetric(sc.totalBackfilledJobsSinceStart, prometheus.GaugeValue, sm.totalBackfilledJobsSinceStart)
	ch <- prometheus.MustNewConstMetric(sc.totalBackfilledJobsSinceCycle, prometheus.GaugeValue, sm.totalBackfilledJobsSinceCycle)
	ch <- prometheus.MustNewConstMetric(sc.totalBackfilledHeterogeneous, prometheus.GaugeValue, sm.totalBackfilledHeterogeneous)

	return nil
}
