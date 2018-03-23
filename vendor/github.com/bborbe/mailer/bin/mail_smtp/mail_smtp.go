package main

import (
	"runtime"

	"fmt"

	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/mailer"
	"github.com/bborbe/mailer/config"
	"github.com/bborbe/mailer/message"
	"github.com/golang/glog"
)

const (
	defaultHost            = "localhost"
	defaultPort            = 1025
	defaultTls             = false
	defaultTlsSkipVerify   = false
	parameterSmtpHost      = "smtp-host"
	parameterSmtpPort      = "smtp-port"
	parameterTls           = "smtp-tls"
	parameterTlsSkipVerify = "smtp-tls-skip-verify"
	parameterFrom          = "from"
	parameterTo            = "to"
	parameterSubject       = "subject"
	parameterBody          = "body"
	parameterAmount        = "amount"
)

var (
	smtpHostPtr          = flag.String(parameterSmtpHost, defaultHost, "smtp host")
	smtpPortPtr          = flag.Int(parameterSmtpPort, defaultPort, "smtp port")
	smtpTlsPtr           = flag.Bool(parameterTls, defaultTls, "smtp tls")
	smtpTlsSkipVerifyPtr = flag.Bool(parameterTlsSkipVerify, defaultTlsSkipVerify, "smtp tls skip verify")
	fromPtr              = flag.String(parameterFrom, "", "from")
	toPtr                = flag.String(parameterTo, "", "to")
	subjectPtr           = flag.String(parameterSubject, "", "subject")
	bodyPtr              = flag.String(parameterBody, "", "body")
	amountPtr            = flag.Int(parameterAmount, 1, "number of mails to send")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := do(); err != nil {
		glog.Exit(err)
	}
}

func do() error {
	mailer, err := createMailer()
	if err != nil {
		return err
	}

	from := *fromPtr
	to := *toPtr
	subject := *subjectPtr
	body := *bodyPtr
	amount := *amountPtr

	message := message.New()
	message.SetSender(from)
	message.SetRecipient(to)
	message.SetSubject(subject)
	message.SetContent(body)

	for i := 0; i < amount; i++ {
		if err := mailer.Send(message); err != nil {
			return err
		}
		glog.V(2).Infof("send %d mail successful", amount)
	}
	return nil
}

func createMailer() (mailer.Mailer, error) {
	smtpHost := *smtpHostPtr
	if len(smtpHost) == 0 {
		return nil, fmt.Errorf("parameter %v missing", parameterSmtpHost)
	}

	smtpPort := *smtpPortPtr
	if smtpPort <= 0 {
		return nil, fmt.Errorf("parameter %v missing", parameterSmtpPort)
	}

	smtpTls := *smtpTlsPtr
	smtpTlsSkipVerify := *smtpTlsSkipVerifyPtr

	config := config.New()
	config.SetSmtpHost(smtpHost)
	config.SetSmtpPort(smtpPort)
	config.SetTls(smtpTls)
	config.SetTlsSkipVerify(smtpTlsSkipVerify)
	return mailer.New(config), nil
}
