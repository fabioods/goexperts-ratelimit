package main

import (
	"github.com/fabioods/goexperts-ratelimit/config"
	"github.com/fabioods/goexperts-ratelimit/internal/pkg/dependencyinjector"
)

func main() {
	configs, err := config.Load(".")
	if err != nil {
		panic(err)
	}

	di := dependencyinjector.NewDependencyInjector(configs)

	deps, err := di.Inject()
	if err != nil {
		panic(err)
	}

	deps.WebServer.Start()
}
