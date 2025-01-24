<a id="markdown-runhttp---prepackaged-runtime-helper-for-http" name="runhttp---prepackaged-runtime-helper-for-http"></a>
# runhttp - Prepackaged Runtime Helper For HTTP
[![GoDoc](https://godoc.org/github.com/asecurityteam/runhttp?status.svg)](https://godoc.org/github.com/asecurityteam/runhttp)

[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=bugs)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=code_smells)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=coverage)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=ncloc)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=alert_status)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=security_rating)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=sqale_index)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_runhttp&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=asecurityteam_runhttp)


<!-- TOC -->

- [runhttp - Prepackaged Runtime Helper For HTTP](#runhttp---prepackaged-runtime-helper-for-http)
    - [Overview](#overview)
    - [Quick Start](#quick-start)
    - [Details](#details)
        - [Configuration](#configuration)
            - [YAML](#yaml)
            - [ENV](#env)
        - [Logging](#logging)
        - [Metrics](#metrics)
    - [Status](#status)
    - [Contributing](#contributing)
        - [Building And Testing](#building-and-testing)
        - [License](#license)
        - [Contributing Agreement](#contributing-agreement)

<!-- /TOC -->

<a id="markdown-overview" name="overview"></a>
## Overview

This project is a tool bundle for running an HTTP service written in go. It comes with
an opinionated choice of logger, metrics client, and configuration parsing. The benefits
are a suite of server metrics built in, a configurable and pluggable shutdown signaling
system, and support for restarting the server without exiting the running process.

<a id="markdown-quick-start" name="quick-start"></a>
## Quick Start

```golang
package main

import (
    "net/http"
    "github.com/asecurityteam/runhttp"
)

func main() {
    // The handler is anything the runtime should serve. This may be a simple
    // handler like this one or any of the many mux/router projects that exist
    // in the go ecosystem.
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        runhttp.LoggerFromContext(r.Context()) // Get a logger
        runtthp.StatFromContext(r.Context()) // Get a metrics client
    })

    // Create any implementation of the settings.Source interface. Here we
    // use the environment variable source.
    source, err := settings.NewEnvSource(os.Environ())
	if err != nil {
		panic(err.Error())
	}

    // Load the runtime using the Source and Handler.
	rt, err := runhttp.New(context.Background(), source, handler)
	if err != nil {
		panic(err.Error())
	}

    // Run the HTTP server.
    if err := rt.Run(); err != nil {
		panic(err.Error())
	}
}
```

<a id="markdown-details" name="details"></a>
## Details

<a id="markdown-configuration" name="configuration"></a>
### Configuration

We use our [settings](https://github.com/asecurityteam/settings) project to manage configuration.
This makes it possible to load configuration values from environment variables, JSON files,
or YAML files. Other configurations sources are possible by implementing the `settings.Source`
interface defined in the settings project.

<a id="markdown-yaml" name="yaml"></a>
#### YAML

If using a YAML file then a configuration would look like:

```yaml
runtime:
  signals:
    # ([]string) Which signal handlers are installed. Choices are OS.
    installed:
      - "OS"
    os:
      # ([]int) Which signals to listen for.
      signals:
        - 15
        - 2
  stats:
    # (string) Destination stream of the stats. One of NULLSTAT, DATADOG.
    output: "DATADOG"
    datadog:
      # (int) Max packet size to send.
      packetsize: 32768
      # ([]string) Any static tags for all metrics.
      tags:
      # (time.Duration) Frequencing of sending metrics to listener.
      flushinterval: "10s"
      # (string) Listener address to use when sending metrics.
      address: "localhost:8125"
  logger:
    # (string) Destination stream of the logs. One of STDOUT, NULL.
    output: "STDOUT"
    # (string) The minimum level of logs to emit. One of DEBUG, INFO, WARN, ERROR.
    level: "INFO"
  connstate:
    # (time.Duration) Interval on which gauges are reported.
    reportinterval: "5s"
    # (string) Name of the counter metric tracking hijacked clients.
    hijackedcounter: "http.server.connstate.hijacked"
    # (string) Name of the counter metric tracking closed clients.
    closedcounter: "http.server.connstate.closed"
    # (string) Name of the gauge metric tracking idle clients.
    idlegauge: "http.server.connstate.idle.gauge"
    # (string) Name of the counter metric tracking idle clients.
    idlecounter: "http.server.connstate.idle"
    # (string) Name of the gauge metric tracking active clients.
    activegauge: "http.server.connstate.active.gauge"
    # (string) Name of the counter metric tracking active clients.
    activecounter: "http.server.connstate.active"
    # (string) Name of the gauge metric tracking new clients.
    newgauge: "http.server.connstate.new.gauge"
    # (string) Name of the counter metric tracking new clients.
    newcounter: "http.server.connstate.new"
  expvar:
    # (string) Name of the metric tracking allocated bytes
    alloc: "go_expvar.memstats.alloc"
    # (string) Name of the metric tracking number of frees
    frees: "go_expvar.memstats.frees"
    # (string) Name of the metric tracking allocated bytes
    heapalloc: "go_expvar.memstats.heap_alloc"
    # (string) Name of the metric tracking bytes in unused spans
    heapidle: "go_expvar.memstats.heap_idle"
    # (string) Name of the metric tracking bytes in in-use spans
    heapinuse: "go_expvar.memstats.heap_inuse"
    # (string) Name of the metric tracking total number of object allocated"
    heapobject: "go_expvar.memstats.heap_objects"
    # (string) Name of the metric tracking bytes realeased to the OS
    heaprealeased: "go_expvar.memstats.heap_released"
    # (string) Name of the metric tracking bytes obtained from the system
    heapstats: "go_expvar.memstats.heap_sys"
    # (string) Name of the metric tracking number of pointer lookups
    lookups: "go_expvar.memstats.lookups"
    # (string) Name of the metric tracking number of mallocs
    mallocs: "go_expvar.memstats.mallocs"
    # (string) Name of the metric tracking number of garbage collections
    numgc: "go_expvar.memstats.num_gc"
    # (string) Name of the metric tracking duration of GC pauses
    pausens: "go_expvar.memstats.pause_ns"
    # (string) Name of the metric tracking total GC pause duration over lifetime process
    pausetotalns: "go_expvar.memstats.pause_total_ns"
    # (string) Name of the metric tracking allocated bytes (even if freed)
    totalalloc: "go_expvar.memstats.total_alloc"
    # (string) Name of the metric tracking number of active go routines
    goroutinesexists: "go_expvar.goroutines.exists"
    # (time.Duration) Interval on which metrics are reported
    reportinterval: "5s"
  httpserver:
    # (string) The listening address of the server.
    address: ":8080"
```

<a id="markdown-env" name="env"></a>
#### ENV

If using the environment variable loader then a configuration would look like:

```bash
# (string) The listening address of the server.
RUNTIME_HTTPSERVER_ADDRESS=":8080"
# (time.Duration) Interval on which gauges are reported.
RUNTIME_CONNSTATE_REPORTINTERVAL="5s"
# (string) Name of the counter metric tracking hijacked clients.
RUNTIME_CONNSTATE_HIJACKEDCOUNTER="http.server.connstate.hijacked"
# (string) Name of the counter metric tracking closed clients.
RUNTIME_CONNSTATE_CLOSEDCOUNTER="http.server.connstate.closed"
# (string) Name of the gauge metric tracking idle clients.
RUNTIME_CONNSTATE_IDLEGAUGE="http.server.connstate.idle.gauge"
# (string) Name of the counter metric tracking idle clients.
RUNTIME_CONNSTATE_IDLECOUNTER="http.server.connstate.idle"
# (string) Name of the gauge metric tracking active clients.
RUNTIME_CONNSTATE_ACTIVEGAUGE="http.server.connstate.active.gauge"
# (string) Name of the counter metric tracking active clients.
RUNTIME_CONNSTATE_ACTIVECOUNTER="http.server.connstate.active"
# (string) Name of the gauge metric tracking new clients.
RUNTIME_CONNSTATE_NEWGAUGE="http.server.connstate.new.gauge"
# (string) Name of the counter metric tracking new clients.
RUNTIME_CONNSTATE_NEWCOUNTER="http.server.connstate.new"
# (string) Name of the metric tracking allocated bytes
RUNTIME_EXPVAR_ALLOC: "go_expvar.memstats.alloc"
# (string) Name of the metric tracking number of frees
RUNTIME_EXPVAR_FREES: "go_expvar.memstats.frees"
# (string) Name of the metric tracking allocated bytes
RUNTIME_EXPVAR_HEAPALLOC: "go_expvar.memstats.heap_alloc"
# (string) Name of the metric tracking bytes in unused spans
RUNTIME_EXPVAR_HEAPIDLE: "go_expvar.memstats.heap_idle"
# (string) Name of the metric tracking bytes in in-use spans
RUNTIME_EXPVAR_HEAPINUSE: "go_expvar.memstats.heap_inuse"
# (string) Name of the metric tracking total number of object allocated"
RUNTIME_EXPVAR_HEAPOBJECT: "go_expvar.memstats.heap_objects"
# (string) Name of the metric tracking bytes realeased to the OS
RUNTIME_EXPVAR_HEAPREALEASED: "go_expvar.memstats.heap_released"
# (string) Name of the metric tracking bytes obtained from the system
RUNTIME_EXPVAR_HEAPSTATS: "go_expvar.memstats.heap_sys"
# (string) Name of the metric tracking number of pointer lookups
RUNTIME_EXPVAR_LOOKUPS: "go_expvar.memstats.lookups"
# (string) Name of the metric tracking number of mallocs
RUNTIME_EXPVAR_MALLOCS: "go_expvar.memstats.mallocs"
# (string) Name of the metric tracking number of garbage collections
RUNTIME_EXPVAR_NUMGC: "go_expvar.memstats.num_gc"
# (string) Name of the metric tracking duration of GC pauses
RUNTIME_EXPVAR_PAUSENS: "go_expvar.memstats.pause_ns"
# (string) Name of the metric tracking total GC pause duration over lifetime process
RUNTIME_EXPVAR_PAUSETOTALNS: "go_expvar.memstats.pause_total_ns"
# (string) Name of the metric tracking allocated bytes (even if freed)
RUNTIME_EXPVAR_TOTALALLOC: "go_expvar.memstats.total_alloc"
# (string) Name of the metric tracking number of active go routines
RUNTIME_EXPVAR_GOROUTINESEXISTS: "go_expvar.goroutines.exists"
# (time.Duration) Interval on which metrics are reported
RUNTIME_EXPVAR_REPORTINTERVAL: "5s"
# (string) Destination stream of the logs. One of STDOUT, NULL.
RUNTIME_LOGGER_OUTPUT="STDOUT"
# (string) The minimum level of logs to emit. One of DEBUG, INFO, WARN, ERROR.
RUNTIME_LOGGER_LEVEL="INFO"
# (string) Destination stream of the stats. One of NULLSTAT, DATADOG.
RUNTIME_STATS_OUTPUT="DATADOG"
# (int) Max packet size to send.
RUNTIME_STATS_DATADOG_PACKETSIZE="32768"
# ([]string) Any static tags for all metrics.
RUNTIME_STATS_DATADOG_TAGS=""
# (time.Duration) Frequencing of sending metrics to listener.
RUNTIME_STATS_DATADOG_FLUSHINTERVAL="10s"
# (string) Listener address to use when sending metrics.
RUNTIME_STATS_DATADOG_ADDRESS="localhost:8125"
# ([]string) Which signal handlers are installed. Choices are OS.
RUNTIME_SIGNALS_INSTALLED="OS"
# ([]int) Which signals to listen for.
RUNTIME_SIGNALS_OS_SIGNALS="15 2"
```

<a id="markdown-logging" name="logging"></a>
### Logging

This project will install a [logevent](https://github.com/asecurityteam/logevent) instance
in the context. From within an HTTP handler the logger should be accessed using
`runhttp.LoggerFromContext(r.Context())`.

<a id="markdown-metrics" name="metrics"></a>
### Metrics

This project will install an [xstats](https://github.com/rs/xstats) instance in the context.
Custom metrics may be emitted by extracting the client using `runhttp.StatFromContext(r.Context())`.

In addition, the server will emit counters for new, active, idle, closed, and hijacked connections.
Note that hijacked in this case refers to a Go behavior defined
[here](https://golang.org/pkg/net/http/#Hijacker). The server emits gauges on an interval for
new, active, and idle connections.

Go runtime metrics are also emitted. These values are extracted on a specified polling interval from the [runtime](https://golang.org/pkg/runtime/#MemStats) package.
The table [here](https://docs.datadoghq.com/integrations/go_expvar/#metrics) illustrates how we expect to see these values as metrics.

<a id="markdown-status" name="status"></a>
## Status

This project is in incubation which means we are not yet operating this tool in production
and the interfaces are subject to change.

<a id="markdown-contributing" name="contributing"></a>
## Contributing

<a id="markdown-building-and-testing" name="building-and-testing"></a>
### Building And Testing

We publish a docker image called [SDCLI](https://github.com/asecurityteam/sdcli) that
bundles all of our build dependencies. It is used by the included Makefile to help make
building and testing a bit easier. The following actions are available through the Makefile:

-   make dep

    Install the project dependencies into a vendor directory

-   make lint

    Run our static analysis suite

-   make test

    Run unit tests and generate a coverage artifact

-   make integration

    Run integration tests and generate a coverage artifact

-   make coverage

    Report the combined coverage for unit and integration tests

<a id="markdown-license" name="license"></a>
### License

This project is licensed under Apache 2.0. See LICENSE.txt for details.

<a id="markdown-contributing-agreement" name="contributing-agreement"></a>
### Contributing Agreement

Atlassian requires signing a contributor's agreement before we can accept a patch. If
you are an individual you can fill out the [individual
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d).
If you are contributing on behalf of your company then please fill out the [corporate
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b).
