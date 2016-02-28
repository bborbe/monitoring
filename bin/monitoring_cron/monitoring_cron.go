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
	"github.com/bborbe/mailer"
	mail_config "github.com/bborbe/mailer/config"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_configuration_parser "github.com/bborbe/monitoring/configuration_parser"
	monitoring_node "github.com/bborbe/monitoring/node"
	monitoring_notifier "github.com/bborbe/monitoring/notifier"
	monitoring_runner_hierarchy "github.com/bborbe/monitoring/runner/hierarchy"
	"github.com/bborbe/webdriver"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_CONFIG   = "config"
)

type Run func(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult

type Notify func(results []monitoring_check.CheckResult) error

type ParseConfiguration func(content []byte) ([]monitoring_node.Node, error)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(log.ERROR), log.FLAG_USAGE)
	configPtr := flag.String(PARAMETER_CONFIG, "", "config")
	smtpUserPtr := flag.String("smtp-user", "smtp@benjamin-borbe.de", "string")
	smtpPasswordPtr := flag.String("smtp-password", "-", "string")
	smtpHostPtr := flag.String("smtp-host", "iredmail.mailfolder.org", "string")
	smtpPortPtr := flag.Int("smtp-port", 465, "int")
	senderPtr := flag.String("sender", "smtp@benjamin-borbe.de", "string")
	recipientPtr := flag.String("recipient", "bborbe@rocketnews.de", "string")
	maxConcurrencyPtr := flag.Int("max", runtime.NumCPU(), "max concurrency")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	logger.Debugf("max concurrency: %d", *maxConcurrencyPtr)

	mailConfig := mail_config.New()
	mailConfig.SetSmtpUser(*smtpUserPtr)
	mailConfig.SetSmtpPassword(*smtpPasswordPtr)
	mailConfig.SetSmtpHost(*smtpHostPtr)
	mailConfig.SetSmtpPort(*smtpPortPtr)
	writer := os.Stdout
	runner := monitoring_runner_hierarchy.New(*maxConcurrencyPtr)
	mailer := mailer.New(mailConfig)
	notifier := monitoring_notifier.New(mailer, *senderPtr, *recipientPtr)
	driver := webdriver.NewPhantomJsDriver("phantomjs")
	driver.Start()
	defer driver.Stop()
	configurationParser := monitoring_configuration_parser.New(driver)

	err := do(writer, runner.Run, notifier.Notify, configurationParser.ParseConfiguration, *configPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(writer io.Writer, run Run, notify Notify, parseConfiguration ParseConfiguration, configPath string) error {
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
	results := make([]monitoring_check.CheckResult, 0)
	failedChecks := 0
	var result monitoring_check.CheckResult
	for result = range run(nodes) {
		if result.Success() {
			fmt.Fprintf(writer, "[OK]   %s (%d ms)\n", result.Message(), result.Duration()/time.Millisecond)
		} else {
			fmt.Fprintf(writer, "[FAIL] %s - %v (%d ms)\n", result.Message(), result.Error(), result.Duration()/time.Millisecond)
			failedChecks++
		}
		results = append(results, result)
	}
	logger.Debugf("all checks executed")
	if failedChecks > 0 {
		fmt.Fprintf(writer, "found %d failed checks => send mail\n", failedChecks)
		err = notify(results)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(writer, "check finished\n")
	return err
}
