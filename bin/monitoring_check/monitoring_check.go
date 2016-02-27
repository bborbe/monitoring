package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	"time"

	io_util "github.com/bborbe/io/util"
	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_configuration_parser "github.com/bborbe/monitoring/configuration_parser"
	monitoring_node "github.com/bborbe/monitoring/node"
	monitoring_runner "github.com/bborbe/monitoring/runner"
	monitoring_runner_all "github.com/bborbe/monitoring/runner/all"
	monitoring_runner_hierarchy "github.com/bborbe/monitoring/runner/hierarchy"
"github.com/bborbe/webdriver"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_CONFIG   = "config"
	PARAMETER_MODE     = "mode"
	PARAMETER_MAX      = "max"
)

type Run func(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult

type ParseConfiguration func(content []byte) ([]monitoring_node.Node, error)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(log.ERROR), log.FLAG_USAGE)
	modePtr := flag.String(PARAMETER_MODE, "", "mode (all|hierachy)")
	configPtr := flag.String(PARAMETER_CONFIG, "", "config")
	maxConcurrencyPtr := flag.Int(PARAMETER_MAX, runtime.NumCPU()*4, "max concurrency")
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
	driver := webdriver.NewPhantomJsDriver("/opt/phantomjs-2.1.1-macosx/bin/phantomjs")
	driver.Start()
	defer driver.Stop()

	configurationParser := monitoring_configuration_parser.New(driver)
	writer := os.Stdout
	defer writer.Close()

	err := do(writer, runner.Run, configurationParser.ParseConfiguration, *configPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(writer io.Writer, run Run, parseConfiguration ParseConfiguration, configPath string) error {
	var err error
	fmt.Fprintf(writer, "check started\n")
	if len(configPath) == 0 {
		return fmt.Errorf("parameter {} missing", PARAMETER_CONFIG)
	}
	path, err := io_util.NormalizePath(configPath)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	nodes, err := parseConfiguration(content)
	if err != nil {
		return err
	}
	var result monitoring_check.CheckResult
	for result = range run(nodes) {
		if result.Success() {
			fmt.Fprintf(writer, "[OK]   %s (%d ms)\n", result.Message(), result.Duration()/time.Millisecond)
		} else {
			fmt.Fprintf(writer, "[FAIL] %s - %v (%d ms)\n", result.Message(), result.Error(), result.Duration()/time.Millisecond)
		}
	}
	fmt.Fprintf(writer, "check finished\n")
	return err
}
