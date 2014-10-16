package notifier

import (
	"bytes"
	"fmt"

	"net/smtp"

	"crypto/tls"
	"net/mail"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/check"
)

var logger = log.DefaultLogger

type MailConfig interface {
	SmtpUser() string
	SmtpPassword() string
	SmtpHost() string
	SmtpPort() int
	Sender() string
	Recipient() string
}

func Notify(mailconfig MailConfig, results []check.CheckResult) error {
	logger.Debug("notify results")
	mailContent := buildMailContent(results)
	err := sendMail(mailconfig, mailContent)
	logger.Debug("mail sent")
	return err
}

func sendMail(mailconfig MailConfig, content string) error {
	logger.Debugf("sendMail to %s", mailconfig.Recipient())
	auth := smtp.PlainAuth(
		"",
		mailconfig.SmtpUser(),
		mailconfig.SmtpPassword(),
		mailconfig.SmtpHost(),
	)
	servername := fmt.Sprintf("%s:%d", mailconfig.SmtpHost(), mailconfig.SmtpPort())
	logger.Debugf("connect to smtp-server to %s", servername)

	from := mail.Address{"", mailconfig.Sender()}
	to := mail.Address{"", mailconfig.Recipient()}
	subj := "Monitoring Result"

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + content

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         servername,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	smtpClient, err := smtp.NewClient(conn, mailconfig.SmtpHost())
	if err != nil {
		return nil
	}

	err = smtpClient.Auth(auth)
	if err != nil {
		return err
	}

	err = smtpClient.Mail(mailconfig.Sender())
	if err != nil {
		return err
	}

	err = smtpClient.Rcpt(mailconfig.Recipient())
	if err != nil {
		return err
	}

	data, err := smtpClient.Data()
	if err != nil {
		return err
	}

	data.Write([]byte(message))

	err = data.Close()
	if err != nil {
		return err
	}

	return smtpClient.Quit()
}

func buildMailContent(results []check.CheckResult) string {
	failures := failures(results)
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Checks executed: %d\n", len(results)))
	buffer.WriteString(fmt.Sprintf("Checks failed: %d\n", len(failures)))
	for _, result := range failures {
		buffer.WriteString(fmt.Sprintf("%s - %v\n", result.Message(), result.Error()))
	}
	logger.Debug("return mailcontent")
	return buffer.String()
}

func failures(results []check.CheckResult) []check.CheckResult {
	failures := make([]check.CheckResult, 0)
	for _, result := range results {
		if !result.Success() {
			failures = append(failures, result)
		}
	}
	return failures
}
