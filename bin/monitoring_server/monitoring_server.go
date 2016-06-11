package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	flag "github.com/bborbe/flagenv"
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
	DEFAULT_LOCK             = "~/.monitoring_cron.lock"
	DEFAULT_DELAY            = time.Minute * 5
	PARAMETER_LOGLEVEL       = "loglevel"
	PARAMETER_CONFIG         = "config"
	PARAMETER_DRIVER         = "driver"
	PARAMETER_DELAY          = "delay"
	PARAMETER_SMTP_USER      = "smtp-user"
	PARAMETER_SMTP_PASSWORD  = "smtp-password"
	PARAMETER_SMTP_HOST      = "smtp-host"
	PARAMETER_SMTP_PORT      = "smtp-port"
	PARAMETER_SMTP_SENDER    = "sender"
	PARAMETER_SMTP_RECIPIENT = "recipient"
	PARAMETER_CONCURRENT     = "max"
	PARAMETER_ONE_TIME       = "one-time"
	PARAMETER_LOCK           = "lock"
)

type Run func(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult

type Notify func(results []monitoring_check.CheckResult) error

type ParseConfiguration func(content []byte) ([]monitoring_node.Node, error)

type ParseNodes func(path string) ([]monitoring_node.Node, error)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(log.ERROR), log.FLAG_USAGE)
	configPtr := flag.String(PARAMETER_CONFIG, "", "config")
	driverPtr := flag.String(PARAMETER_DRIVER, "phantomjs", "driver phantomjs|chromedriver")
	smtpUserPtr := flag.String(PARAMETER_SMTP_USER, "", "string")
	smtpPasswordPtr := flag.String(PARAMETER_SMTP_PASSWORD, "", "string")
	smtpHostPtr := flag.String(PARAMETER_SMTP_HOST, "iredmail.mailfolder.org", "string")
	smtpPortPtr := flag.Int(PARAMETER_SMTP_PORT, 25, "int")
	senderPtr := flag.String(PARAMETER_SMTP_SENDER, "smtp@benjamin-borbe.de", "string")
	recipientPtr := flag.String(PARAMETER_SMTP_RECIPIENT, "bborbe@rocketnews.de", "string")
	maxConcurrencyPtr := flag.Int(PARAMETER_CONCURRENT, runtime.NumCPU(), "max concurrency")
	lockNamePtr := flag.String(PARAMETER_LOCK, DEFAULT_LOCK, "lock file")
	delayPtr := flag.Duration(PARAMETER_DELAY, DEFAULT_DELAY, "delay")
	oneTimePtr := flag.Bool(PARAMETER_ONE_TIME, false, "exit after first backup")

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

	mailConfig := mail_config.New()
	mailConfig.SetSmtpUser(*smtpUserPtr)
	mailConfig.SetSmtpPassword(*smtpPasswordPtr)
	mailConfig.SetSmtpHost(*smtpHostPtr)
	mailConfig.SetSmtpPort(*smtpPortPtr)
	runner := monitoring_runner_hierarchy.New(*maxConcurrencyPtr)
	mailer := mailer.New(mailConfig)
	notifier := monitoring_notifier.New(mailer, *senderPtr, *recipientPtr)
	configurationParser := monitoring_configuration_parser.New(driver)

	err := do(runner.Run, notifier.Notify, func(path string) ([]monitoring_node.Node, error) {
		logger.Debugf("read config")
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		logger.Debugf("parse config")
		nodes, err := configurationParser.ParseConfiguration(content)
		if err != nil {
			return nil, err
		}
		return nodes, nil
	}, *configPtr, *lockNamePtr, *delayPtr, *oneTimePtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(run Run, notify Notify, parseNodes ParseNodes, configPath string, lockName string, delay time.Duration, oneTime bool) error {
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

	if len(configPath) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_CONFIG)
	}
	path, err := io_util.NormalizePath(configPath)
	if err != nil {
		logger.Debugf("normalize path failed: %v", err)
		return err
	}

	for {
		logger.Debugf("check started")

		nodes, err := parseNodes(path)
		if err != nil {
			return fmt.Errorf("parse config failed: %v", err)
		}

		logger.Debugf("run checks")
		results := make([]monitoring_check.CheckResult, 0)
		failedChecks := 0
		var result monitoring_check.CheckResult
		for result = range run(nodes) {
			if !result.Success() {
				failedChecks++
			}
			results = append(results, result)
		}
		logger.Debugf("all checks executed, %d failed", failedChecks)
		if failedChecks > 0 {
			err = notify(results)
			if err != nil {
				return err
			}
		}
		logger.Debugf("check finished")

		if oneTime {
			return nil
		}

		logger.Debugf("sleep for %v", delay)
		time.Sleep(delay)
		logger.Debugf("sleep done")
	}

	return err
}
