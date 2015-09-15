package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bborbe/log"
	"github.com/bborbe/mailer"
	mail_config "github.com/bborbe/mailer/config"
	"github.com/bborbe/monitoring/check"
	monitoring_configuration "github.com/bborbe/monitoring/configuration"
	monitoring_notifier "github.com/bborbe/monitoring/notifier"
	monitoring_runner "github.com/bborbe/monitoring/runner"
	monitoring_runner_hierarchy "github.com/bborbe/monitoring/runner/hierarchy"
)

var logger = log.DefaultLogger

func main() {
	defer logger.Close()
	logLevelPtr := flag.String("loglevel", log.LogLevelToString(log.ERROR), "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	smtpUserPtr := flag.String("smtp-user", "smtp@benjamin-borbe.de", "string")
	smtpPasswordPtr := flag.String("smtp-password", "-", "string")
	smtpHostPtr := flag.String("smtp-host", "iredmail.mailfolder.org", "string")
	smtpPortPtr := flag.Int("smtp-port", 465, "int")
	senderPtr := flag.String("sender", "smtp@benjamin-borbe.de", "string")
	recipientPtr := flag.String("recipient", "bborbe@rocketnews.de", "string")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)
	mailConfig := mail_config.New()
	mailConfig.SetSmtpUser(*smtpUserPtr)
	mailConfig.SetSmtpPassword(*smtpPasswordPtr)
	mailConfig.SetSmtpHost(*smtpHostPtr)
	mailConfig.SetSmtpPort(*smtpPortPtr)
	writer := os.Stdout
	configuration := monitoring_configuration.New()
	runner := monitoring_runner_hierarchy.New()
	mailer := mailer.New(mailConfig)
	notifier := monitoring_notifier.New(mailer, *senderPtr, *recipientPtr)
	err := do(writer, runner, configuration, notifier)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(writer io.Writer, runner monitoring_runner.Runner, configuration monitoring_configuration.Configuration, notifier monitoring_notifier.Notifier) error {
	var err error
	fmt.Fprintf(writer, "check started\n")
	resultChannel := runner.Run(configuration)
	results := make([]check.CheckResult, 0)
	failedChecks := 0
	for result := range resultChannel {
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
