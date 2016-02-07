package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/configuration"
	"github.com/bborbe/monitoring/runner"
	"github.com/bborbe/monitoring/runner/all"
	"github.com/bborbe/monitoring/runner/hierarchy"
	"runtime"
)

var logger = log.DefaultLogger

func main() {
	defer logger.Close()
	logLevelPtr := flag.String("loglevel", log.LogLevelToString(log.ERROR), log.FLAG_USAGE)
	modePtr := flag.String("mode", "", "mode (all|hierachy)")
	maxConcurrencyPtr := flag.Int("max", runtime.NumCPU() * 2, "max concurrency")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	logger.Debugf("max concurrency: %d", *maxConcurrencyPtr)

	var r runner.Runner
	if "all" == *modePtr {
		logger.Debug("runner = all")
		r = all.New(*maxConcurrencyPtr)
	} else {
		logger.Debug("runner = hierarchy")
		r = hierarchy.New(*maxConcurrencyPtr)
	}
	c := configuration.New()
	writer := os.Stdout
	err := do(writer, r, c)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(writer io.Writer, r runner.Runner, c configuration.Configuration) error {
	var err error
	fmt.Fprintf(writer, "check started\n")
	results := r.Run(c)
	for result := range results {
		if result.Success() {
			fmt.Fprintf(writer, "[OK]   %s\n", result.Message())
		} else {
			fmt.Fprintf(writer, "[FAIL] %s - %v\n", result.Message(), result.Error())
		}
	}
	fmt.Fprintf(writer, "check finished\n")
	return err
}
