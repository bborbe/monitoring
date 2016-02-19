package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"runtime"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_configuration "github.com/bborbe/monitoring/configuration"
	monitoring_runner "github.com/bborbe/monitoring/runner"
	monitoring_runner_all "github.com/bborbe/monitoring/runner/all"
	monitoring_runner_hierarchy "github.com/bborbe/monitoring/runner/hierarchy"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL = "loglevel"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(log.ERROR), log.FLAG_USAGE)
	modePtr := flag.String("mode", "", "mode (all|hierachy)")
	maxConcurrencyPtr := flag.Int("max", runtime.NumCPU()*2, "max concurrency")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	logger.Debugf("max concurrency: %d", *maxConcurrencyPtr)

	var r monitoring_runner.Runner
	if "all" == *modePtr {
		logger.Debug("runner = all")
		r = monitoring_runner_all.New(*maxConcurrencyPtr)
	} else {
		logger.Debug("runner = hierarchy")
		r = monitoring_runner_hierarchy.New(*maxConcurrencyPtr)
	}
	c := monitoring_configuration.New()
	writer := os.Stdout
	err := do(writer, r, c)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(writer io.Writer, r monitoring_runner.Runner, c monitoring_configuration.Configuration) error {
	var err error
	fmt.Fprintf(writer, "check started\n")
	results := r.Run(c)
	var result monitoring_check.CheckResult
	for result = range results {
		if result.Success() {
			fmt.Fprintf(writer, "[OK]   %s\n", result.Message())
		} else {
			fmt.Fprintf(writer, "[FAIL] %s - %v\n", result.Message(), result.Error())
		}
	}
	fmt.Fprintf(writer, "check finished\n")
	return err
}
