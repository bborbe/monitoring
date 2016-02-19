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
	monitoring_node "github.com/bborbe/monitoring/node"
	monitoring_runner "github.com/bborbe/monitoring/runner"
	monitoring_runner_all "github.com/bborbe/monitoring/runner/all"
	monitoring_runner_hierarchy "github.com/bborbe/monitoring/runner/hierarchy"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL = "loglevel"
)

type GetNodes func() ([]monitoring_node.Node, error)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(log.ERROR), log.FLAG_USAGE)
	modePtr := flag.String("mode", "", "mode (all|hierachy)")
	maxConcurrencyPtr := flag.Int("max", runtime.NumCPU()*2, "max concurrency")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	logger.Debugf("max concurrency: %d", *maxConcurrencyPtr)

	var runner monitoring_runner.Runner
	if "all" == *modePtr {
		logger.Debug("runner = all")
		runner = monitoring_runner_all.New(*maxConcurrencyPtr)
	} else {
		logger.Debug("runner = hierarchy")
		runner = monitoring_runner_hierarchy.New(*maxConcurrencyPtr)
	}
	configuration := monitoring_configuration.New()
	writer := os.Stdout
	err := do(writer, runner, configuration.Nodes)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(writer io.Writer, runner monitoring_runner.Runner, getNodes GetNodes) error {
	var err error
	fmt.Fprintf(writer, "check started\n")
	nodes, err := getNodes()
	if err != nil {
		return err
	}
	results := runner.Run(nodes)
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
