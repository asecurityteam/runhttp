<a id="markdown-runhttp---prepackaged-runtime-helper-for-http" name="runhttp---prepackaged-runtime-helper-for-http"></a>
# runhttp - Prepackaged Runtime Helper For HTTP

*Status: Incubation*

<!-- TOC -->

- [runhttp - Prepackaged Runtime Helper For HTTP](#runhttp---prepackaged-runtime-helper-for-http)
    - [Overview](#overview)
    - [Quick Start](#quick-start)
    - [Contributing](#contributing)
        - [License](#license)
        - [Contributing Agreement](#contributing-agreement)

<!-- /TOC -->

<a id="markdown-overview" name="overview"></a>
## Overview

This project is a tool bundle for running an HTTP service written in go. It comes with
an opinionated choice of logger, metrics client, and configuration parsing. The benefits
are a suite of server metrics built in, a configurable and pluggable shutdown signaling
system, and support for restarting the server without exiting the running process.

Logging is provided by the `logevent` project and you are highly encouraged to use
the logger provided in the request context rather than another logging library. Metrics are
provided by the `xstats` project and you are highly encouraged to use the client
provided in the request context rather than another metrics library. Configuration is
provided by the `settings` project and you are highly encouraged to use the tooling
it provides to extend any configuration options or add your own additional configurable
components.

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

    // Bind your handler to the Component loader.
    rt := &runhttp.Component{Handler: handler}

    // Create any implementation of the settings.Source interface. Here we
    // use the environment variable source.
    source, err := settings.NewEnvSource(os.Environ())
	if err != nil {
		panic(err.Error())
	}

    // Create a pointer that will contain the runtime.
    runnerDst := new(runhttp.Runner)
    // Load the runtime into the pointer.
	err = settings.NewComponent(ctx, source, rt, runnerDst)
	if err != nil {
		panic(err.Error())
    }

    // Run the HTTP server.
    if err := (*runnerDst).Run(); err != nil {
		panic(err.Error())
	}
}
```

<a id="markdown-contributing" name="contributing"></a>
## Contributing

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
