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
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_configuration_parser "github.com/bborbe/monitoring/config"
	monitoring_node "github.com/bborbe/monitoring/node"
	monitoring_runner "github.com/bborbe/monitoring/runner"
	monitoring_runner_all "github.com/bborbe/monitoring/runner/all"
	monitoring_runner_hierarchy "github.com/bborbe/monitoring/runner/hierarchy"
	"github.com/bborbe/webdriver"
	"github.com/golang/glog"
)

const (
	PARAMETER_CONFIG     = "config"
	PARAMETER_MODE       = "mode"
	PARAMETER_CONCURRENT = "concurrent"
	PARAMETER_DRIVER     = "driver"
)

type Run func(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult

type ParseConfiguration func(content []byte) ([]monitoring_node.Node, error)

var (
	modePtr           = flag.String(PARAMETER_MODE, "", "mode (all|hierachy)")
	configPtr         = flag.String(PARAMETER_CONFIG, "", "config")
	maxConcurrencyPtr = flag.Int(PARAMETER_CONCURRENT, runtime.NumCPU()*4, "max concurrency")
	driverPtr         = flag.String(PARAMETER_DRIVER, "phantomjs", "driver phantomjs|chromedriver")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()

	glog.V(2).Infof("max concurrency: %d", *maxConcurrencyPtr)

	var driver webdriver.WebDriver
	if *driverPtr == "chromedriver" {
		driver = webdriver.NewChromeDriver("chromedriver")
	} else {
		driver = webdriver.NewPhantomJsDriver("phantomjs")
	}
	driver.Start()
	defer driver.Stop()

	writer := os.Stdout
	var runner monitoring_runner.Runner
	if "all" == *modePtr {
		glog.V(2).Info("runner = all")
		runner = monitoring_runner_all.New(*maxConcurrencyPtr)
	} else {
		glog.V(2).Info("runner = hierarchy")
		runner = monitoring_runner_hierarchy.New(*maxConcurrencyPtr)
	}
	configurationParser := monitoring_configuration_parser.New(driver)

	err := do(writer, runner.Run, configurationParser.ParseConfiguration, *configPtr)
	if err != nil {
		glog.Exit(err)
	}
	glog.V(2).Info("done")
}

func do(writer io.Writer, run Run, parseConfiguration ParseConfiguration, configPath string) error {
	var err error
	start := time.Now()
	fmt.Fprintf(writer, "checks started\n")
	if len(configPath) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_CONFIG)
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
	var success int
	var total int
	for result = range run(nodes) {
		total++
		if result.Success() {
			fmt.Fprintf(writer, "[OK]   %s (%dms)\n", result.Message(), result.Duration()/time.Millisecond)
			success++
		} else {
			fmt.Fprintf(writer, "[FAIL] %s - %v (%dms)\n", result.Message(), result.Error(), result.Duration()/time.Millisecond)
		}
	}
	duration := time.Now().Sub(start) / time.Millisecond
	fmt.Fprintf(writer, "checks finished with %d/%d successful (%dms)\n", success, total, duration)
	return err
}
