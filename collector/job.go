/*
	Copyright 2023 Atrestis Karalis
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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type JobIdMetrics struct {
	JobID   string
	JobName string
	User    string
	Elapsed float64
}

func ParseJobLine(line string) (*JobIdMetrics, error) {
	fields := strings.Fields(line)
	if len(fields) >= 4 {
		elapsed, err := ParseElapsedTime(fields[3])
		if err != nil {
			return nil, fmt.Errorf("error parsing job: %w", err)
		}

		return &JobIdMetrics{
			JobID:   fields[0],
			JobName: fields[1],
			User:    fields[2],
			Elapsed: elapsed,
		}, nil
	}
	return nil, fmt.Errorf("no job lines found")
}

func ParseJobMetrics(logger log.Logger, input []byte) []*JobIdMetrics {
	var metrics []*JobIdMetrics

	for _, line := range SplitLines(input) {
		jobMetrics, err := ParseJobLine(line)
		if jobMetrics != nil {
			metrics = append(metrics, jobMetrics)
		} else if err == nil {
			level.Warn(logger).Log("msg", "Unable to parser job line", "err", err)
		}
	}

	return metrics
}

func ParseElapsedTime(elapsedStr string) (float64, error) {
	parts := strings.Split(elapsedStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid elapsed time format: %s", elapsedStr)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	elapsedSeconds := float64(hours*3600 + minutes*60 + seconds)
	return elapsedSeconds, nil
}

type JobCollector struct {
	jobInfo *prometheus.Desc
	logger  log.Logger
}

func init() {
	registerCollector("job", defaultDisabled, NewJobCollector)
}

func NewJobCollector(logger log.Logger) (Collector, error) {
	return &JobCollector{
		logger:  logger,
		jobInfo: prometheus.NewDesc("slurm_job_info", "Slurm Job Information", []string{"JobID", "JobName", "User"}, nil),
	}, nil
}

func (jc *JobCollector) Collect(ch chan<- prometheus.Metric) error {
	// Calculate the time one hour ago
	oneHourAgoTime := time.Now().Add(-30 * time.Hour)
	currentTime := time.Now()

	out, err := RunCommand("sacct", "--state=COMPLETED",
		"-S"+oneHourAgoTime.Format("2006-01-02T15:04:05"),
		"-E"+currentTime.Format("2006-01-02T15:04:05"),
		"-X", "-n", "-a",
		"--format=JobID,JobName,User,Elapsed")
	if err != nil {
		return err
	}

	jobMetrics := ParseJobMetrics(jc.logger, out)
	for _, metric := range jobMetrics {
		if metric == nil {
			level.Warn(jc.logger).Log("msg", "Skipping nil metric")
			continue
		}
		ch <- prometheus.MustNewConstMetric(jc.jobInfo, prometheus.GaugeValue, metric.Elapsed, metric.JobID, metric.JobName, metric.User)

		level.Info(jc.logger).Log("msg", fmt.Sprintf("Exported metrics for JobID: %s, JobName: %s, User: %s, Elapsed: %f", metric.JobID, metric.JobName, metric.User, metric.Elapsed))
	}

	return nil
}
