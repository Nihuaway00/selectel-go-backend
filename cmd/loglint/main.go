package main

import (
	"example.com/loglint/pkg/loglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	plugin, err := loglint.New(nil)
	if err != nil {
		panic(err)
	}

	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		panic(err)
	}

	singlechecker.Main(analyzers[0])
}
