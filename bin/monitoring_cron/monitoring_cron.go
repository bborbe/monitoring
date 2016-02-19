package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"runtime"

	"io/ioutil"

	io_util "github.com/bborbe/io/util"
	"github.com/bborbe/log"
	"github.com/bborbe/mailer"
	mail_config "github.com/bborbe/mailer/config"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_configuration "github.com/bborbe/monitoring/configuration"
	"github.com/bborbe/monitoring/configuration_parser"
	monitoring_node "github.com/bborbe/monitoring/node"
	monitoring_notifier "github.com/bborbe/monitoring/notifier"
	monitoring_runner "github.com/bborbe/monitoring/runner"
	monitoring_runner_hierarchy "github.com/bborbe/monitoring/runner/hierarchy"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_CONFIG   = "config"
)

type GetNodes func() ([]monitoring_node.Node, error)

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
	maxConcurrencyPtr := flag.Int("max", runtime.NumCPU()*2, "max concurrency")
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
	var getNodes GetNodes
	if len(*configPtr) > 0 {
		configurationParser := configuration_parser.New()
		getNodes = func() ([]monitoring_node.Node, error) {
			path, err := io_util.NormalizePath(*configPtr)
			if err != nil {
				return nil, err
			}
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}
			return configurationParser.ParseConfiguration(content)
		}
	} else {
		configuration := monitoring_configuration.New()
		getNodes = configuration.Nodes
	}

	err := do(writer, runner, getNodes, notifier)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(writer io.Writer, runner monitoring_runner.Runner, getNodes GetNodes, notifier monitoring_notifier.Notifier) error {
	var err error
	fmt.Fprintf(writer, "check started\n")
	nodes, err := getNodes()
	if err != nil {
		return err
	}
	resultChannel := runner.Run(nodes)
	results := make([]monitoring_check.CheckResult, 0)
	failedChecks := 0
	var result monitoring_check.CheckResult
	for result = range resultChannel {
		if result.Success() {
			fmt.Fprintf(writer, "[OK]   %s\n", result.Message())
		} else {
			fmt.Fprintf(writer, "[FAIL] %s - %v\n", result.Message(), result.Error())
			failedChecks++
		}
		results = append(results, result)
	}
	logger.Debugf("all checks executed")
	if failedChecks > 0 {
		fmt.Fprintf(writer, "found %d failed checks => send mail\n", failedChecks)
		err = notifier.Notify(results)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(writer, "check finished\n")
	return err
}
