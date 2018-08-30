package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"time"

	"context"

	"github.com/bborbe/cron"
	flag "github.com/bborbe/flagenv"
	io_util "github.com/bborbe/io/util"
	"github.com/bborbe/lock"
	"github.com/bborbe/mailer"
	mail_config "github.com/bborbe/mailer/config"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_configuration_parser "github.com/bborbe/monitoring/config"
	monitoring_node "github.com/bborbe/monitoring/node"
	monitoring_notifier "github.com/bborbe/monitoring/notifier"
	monitoring_runner_hierarchy "github.com/bborbe/monitoring/runner/hierarchy"
	"github.com/bborbe/webdriver"
	"github.com/golang/glog"
)

const (
	DEFAULT_LOCK                   = "~/.monitoring_cron.lock"
	DEFAULT_DELAY                  = time.Minute * 5
	PARAMETER_CONFIG               = "config"
	PARAMETER_DRIVER               = "driver"
	PARAMETER_DELAY                = "delay"
	PARAMETER_SMTP_USER            = "smtp-user"
	PARAMETER_SMTP_PASSWORD        = "smtp-password"
	PARAMETER_SMTP_HOST            = "smtp-host"
	PARAMETER_SMTP_PORT            = "smtp-port"
	PARAMETER_SMTP_SENDER          = "sender"
	PARAMETER_SMTP_RECIPIENT       = "recipient"
	PARAMETER_CONCURRENT           = "concurrent"
	PARAMETER_ONE_TIME             = "one-time"
	PARAMETER_LOCK                 = "lock"
	PARAMETER_SMTP_TLS             = "smtp-tls"
	PARAMETER_SMTP_TLS_SKIP_VERIFY = "smtp-tls-skip-verify"
	PARAMETER_SUBJECT              = "subject"
	parameterCronExpression        = "expression"
)

type Run func(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult

type Notify func(sender string, recipient string, subject string, results []monitoring_check.CheckResult) error

type ParseNodes func(path string) ([]monitoring_node.Node, error)

var (
	configPtr         = flag.String(PARAMETER_CONFIG, "", "config")
	driverPtr         = flag.String(PARAMETER_DRIVER, "phantomjs", "driver phantomjs|chromedriver")
	smtpUserPtr       = flag.String(PARAMETER_SMTP_USER, "", "string")
	smtpPasswordPtr   = flag.String(PARAMETER_SMTP_PASSWORD, "", "string")
	smtpHostPtr       = flag.String(PARAMETER_SMTP_HOST, "mail.benjamin-borbe.de", "string")
	smtpPortPtr       = flag.Int(PARAMETER_SMTP_PORT, 25, "int")
	senderPtr         = flag.String(PARAMETER_SMTP_SENDER, "smtp@benjamin-borbe.de", "string")
	recipientPtr      = flag.String(PARAMETER_SMTP_RECIPIENT, "bborbe@rocketnews.de", "string")
	maxConcurrencyPtr = flag.Int(PARAMETER_CONCURRENT, runtime.NumCPU(), "max concurrency")
	lockNamePtr       = flag.String(PARAMETER_LOCK, DEFAULT_LOCK, "lock file")
	cronDelayPtr      = flag.Duration(PARAMETER_DELAY, DEFAULT_DELAY, "delay")
	cronExpressionPtr = flag.String(parameterCronExpression, "", "cron expression '* * * * * ?'")
	cronOneTimePtr    = flag.Bool(PARAMETER_ONE_TIME, false, "exit after first backup")
	tlsPtr            = flag.Bool(PARAMETER_SMTP_TLS, false, "tls")
	tlsSkipVerifyPtr  = flag.Bool(PARAMETER_SMTP_TLS_SKIP_VERIFY, false, "tls skip verify")
	subjectPtr        = flag.String(PARAMETER_SUBJECT, "Monitoring Result", "subject")
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

	runner := monitoring_runner_hierarchy.New(*maxConcurrencyPtr)
	mailer := createMailer()
	notifier := monitoring_notifier.New(mailer)
	configurationParser := monitoring_configuration_parser.New(driver)

	s := &server{
		runner: runner.Run,
		notify: notifier.Notify,
		parseNodes: func(path string) ([]monitoring_node.Node, error) {
			glog.V(2).Infof("read config")
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}
			glog.V(2).Infof("parse config")
			nodes, err := configurationParser.ParseConfiguration(content)
			if err != nil {
				return nil, err
			}
			return nodes, nil
		},
		configPath: *configPtr,
		lockName:   *lockNamePtr,
		delay:      *cronDelayPtr,
		expression: *cronExpressionPtr,
		oneTime:    *cronOneTimePtr,
		sender:     *senderPtr,
		recipient:  *recipientPtr,
		subject:    *subjectPtr,
	}

	if err := s.validate(); err != nil {
		glog.Exit(err)
	}

	if err := s.run(context.Background()); err != nil {
		glog.Exit(err)
	}
	glog.V(2).Info("done")
}

type server struct {
	runner     Run
	notify     Notify
	parseNodes ParseNodes
	configPath string
	lockName   string
	delay      time.Duration
	expression string
	oneTime    bool
	sender     string
	recipient  string
	subject    string
}

func createMailer() mailer.Mailer {
	mailConfig := mail_config.New()
	mailConfig.SetSmtpUser(*smtpUserPtr)
	mailConfig.SetSmtpPassword(*smtpPasswordPtr)
	mailConfig.SetSmtpHost(*smtpHostPtr)
	mailConfig.SetSmtpPort(*smtpPortPtr)
	mailConfig.SetTls(*tlsPtr)
	mailConfig.SetTlsSkipVerify(*tlsSkipVerifyPtr)
	return mailer.New(mailConfig)
}

func (s *server) run(ctx context.Context) error {
	glog.V(1).Infof("monitoring server started")

	var err error
	s.lockName, err = io_util.NormalizePath(s.lockName)
	if err != nil {
		return err
	}
	glog.V(2).Infof("try locking %s", s.lockName)
	l := lock.NewLock(s.lockName)
	if err = l.Lock(); err != nil {
		glog.V(2).Infof("lock %s failed: %v", s.lockName, err)
		return err
	}
	defer func() {
		if err := l.Unlock(); err != nil {
			glog.Warningf("unlock failed: %v", err)
		}
	}()

	return cron.NewCronJob(s.oneTime, s.expression, s.delay, s.action).Run(ctx)
}

func (s *server) validate() error {
	if len(s.configPath) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_CONFIG)
	}

	return nil
}

func (s *server) action(ctx context.Context) error {
	glog.V(2).Infof("check started")
	path, err := io_util.NormalizePath(s.configPath)
	if err != nil {
		glog.V(2).Infof("normalize path failed: %v", err)
		return err
	}

	nodes, err := s.parseNodes(path)
	if err != nil {
		return fmt.Errorf("parse config failed: %v", err)
	}

	glog.V(2).Infof("run checks")
	results := make([]monitoring_check.CheckResult, 0)
	var failedChecks []string
	var result monitoring_check.CheckResult
	for result = range s.runner(nodes) {
		if !result.Success() {
			failedChecks = append(failedChecks, result.Message())
		}
		results = append(results, result)
	}
	glog.V(1).Infof("all checks executed, %d failed: %v", len(failedChecks), failedChecks)
	if len(failedChecks) > 0 {
		err = s.notify(s.sender, s.recipient, s.subject, results)
		if err != nil {
			return err
		}
	}
	glog.V(2).Infof("check finished")
	return nil
}
