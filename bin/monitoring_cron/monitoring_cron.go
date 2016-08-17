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
	"github.com/bborbe/lock"
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
	PARAMETER_LOGLEVEL             = "loglevel"
	PARAMETER_CONFIG               = "config"
	PARAMETER_DRIVER               = "driver"
	DEFAULT_LOCK                   = "~/.monitoring_cron.lock"
	PARAMETER_SMTP_USER            = "smtp-user"
	PARAMETER_SMTP_PASSWORD        = "smtp-password"
	PARAMETER_SMTP_HOST            = "smtp-host"
	PARAMETER_SMTP_PORT            = "smtp-port"
	PARAMETER_SMTP_SENDER          = "sender"
	PARAMETER_SMTP_RECIPIENT       = "recipient"
	PARAMETER_CONCURRENT           = "concurrent"
	PARAMETER_LOCK                 = "lock"
	PARAMETER_SMTP_TLS             = "smtp-tls"
	PARAMETER_SMTP_TLS_SKIP_VERIFY = "smtp-tls-skip-verify"
	PARAMETER_SUBJECT              = "subject"
)

type Run func(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult

type Notify func(sender string, recipient string, subject string, results []monitoring_check.CheckResult) error

type ParseConfiguration func(content []byte) ([]monitoring_node.Node, error)

var (
	logLevelPtr       = flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(log.ERROR), log.FLAG_USAGE)
	configPtr         = flag.String(PARAMETER_CONFIG, "", "config")
	driverPtr         = flag.String(PARAMETER_DRIVER, "phantomjs", "driver phantomjs|chromedriver")
	smtpUserPtr       = flag.String(PARAMETER_SMTP_USER, "smtp@benjamin-borbe.de", "string")
	smtpPasswordPtr   = flag.String(PARAMETER_SMTP_PASSWORD, "-", "string")
	smtpHostPtr       = flag.String(PARAMETER_SMTP_HOST, "iredmail.mailfolder.org", "string")
	smtpPortPtr       = flag.Int(PARAMETER_SMTP_PORT, 465, "int")
	senderPtr         = flag.String(PARAMETER_SMTP_SENDER, "smtp@benjamin-borbe.de", "string")
	recipientPtr      = flag.String(PARAMETER_SMTP_RECIPIENT, "bborbe@rocketnews.de", "string")
	maxConcurrencyPtr = flag.Int(PARAMETER_CONCURRENT, runtime.NumCPU(), "max concurrency")
	lockNamePtr       = flag.String(PARAMETER_LOCK, DEFAULT_LOCK, "lock file")
	tlsPtr            = flag.Bool(PARAMETER_SMTP_TLS, false, "tls")
	tlsSkipVerifyPtr  = flag.Bool(PARAMETER_SMTP_TLS_SKIP_VERIFY, false, "tls skip verify")
	subjectPtr        = flag.String(PARAMETER_SUBJECT, "Monitoring Result", "subject")
)

func main() {
	defer logger.Close()
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	logger.Debugf("max concurrency: %d", *maxConcurrencyPtr)

	var driver webdriver.WebDriver
	if *driverPtr == "chromedriver" {
		driver = webdriver.NewChromeDriver("chromedriver")
	} else {
		driver = webdriver.NewPhantomJsDriver("phantomjs")
	}
	driver.Start()
	defer driver.Stop()

	writer := os.Stdout
	mailConfig := mail_config.New()
	mailConfig.SetSmtpUser(*smtpUserPtr)
	mailConfig.SetSmtpPassword(*smtpPasswordPtr)
	mailConfig.SetSmtpHost(*smtpHostPtr)
	mailConfig.SetSmtpPort(*smtpPortPtr)
	mailConfig.SetTls(*tlsPtr)
	mailConfig.SetTlsSkipVerify(*tlsSkipVerifyPtr)

	runner := monitoring_runner_hierarchy.New(*maxConcurrencyPtr)
	mailer := mailer.New(mailConfig)
	notifier := monitoring_notifier.New(mailer)
	configurationParser := monitoring_configuration_parser.New(driver)

	err := do(writer, runner.Run, notifier.Notify, configurationParser.ParseConfiguration, *configPtr, *lockNamePtr, *senderPtr, *recipientPtr, *subjectPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(
	writer io.Writer,
	run Run,
	notify Notify,
	parseConfiguration ParseConfiguration,
	configPath string,
	lockName string,
	sender string,
	recipient string,
	subject string,
) error {
	var err error
	lockName, err = io_util.NormalizePath(lockName)
	if err != nil {
		return err
	}
	logger.Debugf("try locking %s", lockName)
	l := lock.NewLock(lockName)
	if err = l.Lock(); err != nil {
		logger.Debugf("lock %s failed: %v", lockName, err)
		return err
	}
	defer l.Unlock()

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
			fmt.Fprintf(writer, "[OK]   %s (%dms)\n", result.Message(), result.Duration()/time.Millisecond)
		} else {
			fmt.Fprintf(writer, "[FAIL] %s - %v (%dms)\n", result.Message(), result.Error(), result.Duration()/time.Millisecond)
			failedChecks++
		}
		results = append(results, result)
	}
	logger.Debugf("all checks executed")
	if failedChecks > 0 {
		fmt.Fprintf(writer, "found %d failed checks => send mail\n", failedChecks)
		err = notify(sender, recipient, subject, results)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(writer, "check finished\n")
	return err
}
