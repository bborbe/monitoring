package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/configuration"
	"github.com/bborbe/monitoring/notifier"
	"github.com/bborbe/monitoring/runner/all"
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
	mailConfig := new(mailConfig)
	mailConfig.smtpUser = *smtpUserPtr
	mailConfig.smtpPassword = *smtpPasswordPtr
	mailConfig.smtpHost = *smtpHostPtr
	mailConfig.smtpPort = *smtpPortPtr
	mailConfig.sender = *senderPtr
	mailConfig.recipient = *recipientPtr
	writer := os.Stdout
	c := configuration.New()
	err := do(writer, c, mailConfig)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(writer io.Writer, cfg configuration.Configuration, mailConfig notifier.MailConfig) error {
	var err error
	fmt.Fprintf(writer, "check started\n")
	runner := all.New()
	resultChannel := runner.Run(cfg)
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
		err = notifier.Notify(mailConfig, results)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(writer, "check finished\n")
	return err
}

type mailConfig struct {
	smtpUser     string
	smtpPassword string
	smtpHost     string
	smtpPort     int
	sender       string
	recipient    string
}

func (c *mailConfig) SmtpUser() string     { return c.smtpUser }
func (c *mailConfig) SmtpPassword() string { return c.smtpPassword }
func (c *mailConfig) SmtpHost() string     { return c.smtpHost }
func (c *mailConfig) SmtpPort() int        { return c.smtpPort }
func (c *mailConfig) Sender() string       { return c.sender }
func (c *mailConfig) Recipient() string    { return c.recipient }
