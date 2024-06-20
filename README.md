# Prometheus Slurm Exporter

Prometheus collector and exporter for metrics extracted from the [Slurm](https://slurm.schedmd.com/overview.html) resource scheduling system.

## Changes in this fork
* Merged pending Pull Requests from original repository ([#43], [#53], [#54], [#62], [#65], [#70], [#83], [#94], [#96], [#99], [#106]);
* Updated build tools and dependencies following `node_exporter` practices;
* Each collector can be individually enabled/disabled using command line options;
* Refactored code to improve readability and maintainability;
* Refined and ensured fixtures and tests for each collector;
* Automatic test and build using GitHub Actions;
* Changed default port to `9341` (following [default port allocations](https://github.com/prometheus/prometheus/wiki/Default-port-allocations));
* Targeted minimum supported Slurm version 18.08;

[#106]: https://github.com/vpenso/prometheus-slurm-exporter/pull/106
[#99]: https://github.com/vpenso/prometheus-slurm-exporter/pull/99
[#96]: https://github.com/vpenso/prometheus-slurm-exporter/pull/96
[#94]: https://github.com/vpenso/prometheus-slurm-exporter/pull/94
[#83]: https://github.com/vpenso/prometheus-slurm-exporter/pull/83
[#70]: https://github.com/vpenso/prometheus-slurm-exporter/pull/70
[#65]: https://github.com/vpenso/prometheus-slurm-exporter/pull/65
[#62]: https://github.com/vpenso/prometheus-slurm-exporter/pull/62
[#54]: https://github.com/vpenso/prometheus-slurm-exporter/pull/54
[#53]: https://github.com/vpenso/prometheus-slurm-exporter/pull/53
[#43]: https://github.com/vpenso/prometheus-slurm-exporter/pull/43

## Building from source

You need a Go development environment. Then, simply run `make` to build the executables:

    make build

## Installation

* Download executable from [GitHub Releases](https://github.com/astappiev/slurm_exporter/releases) or build from sources; 
* Copy the executable `slurm_exporter` to a node with access to the Slurm command-line interface;
* Check options to disable/enable collectors and set the port to listen to `slurm_exporter -h`;
* A [Systemd Unit](https://www.freedesktop.org/software/systemd/man/systemd.service.html) file to run the executable as service is available in [examples/systemd/slurm_exporter.service](examples/systemd/slurm_exporter.service).

## Exported Metrics

### State of the CPUs

* **Allocated**: CPUs which have been allocated to a job.
* **Idle**: CPUs not allocated to a job and thus available for use.
* **Other**: CPUs which are unavailable for use at the moment.
* **Total**: total number of CPUs.

- Information extracted from the SLURM [**sinfo**](https://slurm.schedmd.com/sinfo.html) command.
- [Slurm CPU Management User and Administrator Guide](https://slurm.schedmd.com/cpu_management.html)

### State of the GPUs

* **Allocated**: GPUs which have been allocated to a job.
* **Other**: GPUs which are unavailable for use at the moment.
* **Total**: total number of GPUs.
* **Utilization**: total GPU utilization on the cluster.

- Information extracted from the SLURM [**sinfo**](https://slurm.schedmd.com/sinfo.html) and [**sacct**](https://slurm.schedmd.com/sacct.html) command.
- [Slurm GRES scheduling](https://slurm.schedmd.com/gres.html)

**NOTE**: The collectors are managed in similar way to `node_exporter`, to disable the GPU collector, use the following command line option `--no-collector.gpus`.

### State of the Nodes

* **Allocated**: nodes which has been allocated to one or more jobs.
* **Completing**: all jobs associated with these nodes are in the process of being completed.
* **Down**: nodes which are unavailable for use.
* **Drain**: with this metric two different states are accounted for:
  - nodes in ``drained`` state (marked unavailable for use per system administrator request)
  - nodes in ``draining`` state (currently executing jobs but which will not be allocated for new ones).
* **Fail**: these nodes are expected to fail soon and are unavailable for use per system administrator request.
* **Error**: nodes which are currently in an error state and not capable of running any jobs.
* **Idle**: nodes not allocated to any jobs and thus available for use.
* **Maint**: nodes which are currently marked with the __maintenance__ flag.
* **Mixed**: nodes which have some of their CPUs ALLOCATED while others are IDLE.
* **Resv**: these nodes are in an advanced reservation and not generally available.

- Information extracted from the SLURM [**sinfo**](https://slurm.schedmd.com/sinfo.html) command.

#### Additional info about node usage

* CPUs: how many are _allocated_, _idle_, _other_ and in _total_.
* Memory: _allocated_ and in _total_.
* Labels: hostname and its Slurm status (e.g. _idle_, _mix_, _allocated_, _draining_, etc.).

### Status of the Jobs

* **PENDING**: Jobs awaiting for resource allocation.
* **PENDING_DEPENDENCY**: Jobs awaiting because of an unexecuted job dependency.
* **RUNNING**: Jobs currently allocated.
* **SUSPENDED**: Job has an allocation but execution has been suspended and CPUs have been released for other jobs.
* **CANCELLED**: Jobs which were explicitly cancelled by the user or system administrator.
* **COMPLETING**: Jobs which are in the process of being completed.
* **COMPLETED**: Jobs have terminated all processes on all nodes with an exit code of zero.
* **CONFIGURING**: Jobs have been allocated resources, but are waiting for them to become ready for use.
* **FAILED**: Jobs terminated with a non-zero exit code or other failure condition.
* **TIMEOUT**: Jobs terminated upon reaching their time limit.
* **PREEMPTED**: Jobs terminated due to preemption.
* **NODE_FAIL**: Jobs terminated due to failure of one or more allocated nodes.

- Information extracted from the SLURM [**squeue**](https://slurm.schedmd.com/squeue.html) command.

### State of the Partitions

* Running/suspended Jobs per partitions, divided between Slurm accounts and users.
* CPUs total/allocated/idle per partition plus used CPU per user ID.

### Jobs information per Account and User

The following information about jobs are also extracted via [squeue](https://slurm.schedmd.com/squeue.html):

* **Running/Pending/Suspended** jobs per SLURM Account.
* **Running/Pending** CPUs per SLURM Account.
* **Running/Pending/Suspended** jobs per SLURM User.
* **Running/Pending** CPUs per SLURM User.

### Scheduler Information

* **Server Thread count**: The number of current active ``slurmctld`` threads.
* **Queue size**: The length of the scheduler queue.
* **DBD Agent queue size**: The length of the message queue for _SlurmDBD_.
* **Last cycle**: Time in microseconds for last scheduling cycle.
* **Mean cycle**: Mean of scheduling cycles since last reset.
* **Cycles per minute**: Counter of scheduling executions per minute.
* **(Backfill) Last cycle**: Time in microseconds of last backfilling cycle.
* **(Backfill) Mean cycle**: Mean of backfilling scheduling cycles in microseconds since last reset.
* **(Backfill) Depth mean**: Mean of processed jobs during backfilling scheduling cycles since last reset.
* **(Backfill) Total Backfilled Jobs** (since last slurm start): number of jobs started thanks to backfilling since last Slurm start.
* **(Backfill) Total Backfilled Jobs** (since last stats cycle start): number of jobs started thanks to backfilling since last time stats where reset.
* **(Backfill) Total backfilled heterogeneous Job components**: number of heterogeneous job components started thanks to backfilling since last Slurm start.

- Information extracted from the SLURM [**sdiag**](https://slurm.schedmd.com/sdiag.html) command.

*DBD Agent queue size*: it is particularly important to keep track of it, since an increasing number of messages
counted with this parameter almost always indicates three issues:
* the _SlurmDBD_ daemon is down;
* the database is either down or unreachable;
* the status of the Slurm accounting DB may be inconsistent (e.g. ``sreport`` missing data, weird utilization of the cluster, etc.).

### Share Information

Collect _share_ statistics for every Slurm account. Refer to the [manpage of the sshare command](https://slurm.schedmd.com/sshare.html) to get more information.

## Prometheus Configuration for the SLURM exporter

It is strongly advisable to configure the Prometheus server with the following parameters:

```
scrape_configs:

#
# SLURM resource manager:
#
  - job_name: 'slurm_exporter'
    scrape_interval:  30s
    scrape_timeout:   30s
    static_configs:
      - targets: ['slurm_host.fqdn:9341']
```

* **scrape_interval**: a 30 seconds interval will avoid possible 'overloading' on the SLURM master due to frequent calls of sdiag/squeue/sinfo commands through the exporter.
* **scrape_timeout**: on a busy SLURM master a too short scraping timeout will abort the communication from the Prometheus server toward the exporter, thus generating a ``context_deadline_exceeded`` error.

The previous configuration file can be immediately used with a fresh installation of Prometheus. At the same time, we highly recommend to include at least the ``global`` section into the configuration. Official documentation about __configuring Prometheus__ is [available here](https://prometheus.io/docs/prometheus/latest/configuration/configuration/).

**NOTE**: the Prometheus server is using __YAML__ as format for its configuration file, thus **indentation** is really important. Before reloading the Prometheus server it would be better to check the syntax:

```
$~ promtool check-config prometheus.yml

Checking prometheus.yml
  SUCCESS: 1 rule files found
[...]
```

## Grafana Dashboard

A [dashboard](https://grafana.com/dashboards/4323) is available in order to
visualize the exported metrics through [Grafana](https://grafana.com):

![Status of the Nodes](https://github.com/vpenso/prometheus-slurm-exporter/raw/master/images/Node_Status.png)

![Status of the Jobs](https://github.com/vpenso/prometheus-slurm-exporter/raw/master/images/Job_Status.png)

![SLURM Scheduler Information](https://github.com/vpenso/prometheus-slurm-exporter/raw/master/images/Scheduler_Info.png)


## License

Copyright 2017-2024 Victor Penso, Matteo Dessalvi, Oleh Astappiev

This is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see http://www.gnu.org/licenses/.
