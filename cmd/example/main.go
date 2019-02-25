package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/asecurityteam/runhttp"
	"github.com/asecurityteam/settings"
)

func main() {
	ctx := context.Background()
	runnerDst := new(runhttp.Runner)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rt := &runhttp.Component{Handler: handler}

	fs := flag.NewFlagSet("example", flag.ContinueOnError)
	fs.Usage = func() {}
	err := fs.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		rg, _ := settings.Convert(rt.Settings())
		fmt.Println(settings.ExampleEnvGroups([]settings.Group{rg}))
		return
	}

	source, err := settings.NewEnvSource(os.Environ())
	if err != nil {
		panic(err.Error())
	}

	err = settings.NewComponent(ctx, source, rt, runnerDst)
	if err != nil {
		panic(err.Error())
	}

	if err := (*runnerDst).Run(); err != nil {
		panic(err.Error())
	}
}
