package notifier

import (
	"bytes"
	"fmt"

	"github.com/bborbe/log"
	"github.com/bborbe/mailer"
	"github.com/bborbe/mailer/message"
	monitoring_check "github.com/bborbe/monitoring/check"
)

var logger = log.DefaultLogger

type notifier struct {
	mailer mailer.Mailer
}

type Notifier interface {
	Notify(results []monitoring_check.CheckResult) error
}

func New(mailer mailer.Mailer) *notifier {
	n := new(notifier)
	n.mailer = mailer
	return n
}

func (n *notifier) Notify(sender string, recipient string, subject string, results []monitoring_check.CheckResult) error {
	logger.Debug("notify results")
	mailContent := buildMailContent(results)
	message := buildMessage(sender, recipient, subject, mailContent)
	err := n.mailer.Send(message)
	logger.Debug("mail sent")
	return err
}

func buildMessage(sender string, recipient string, subject, content string) mailer.Message {
	m := message.New()
	m.SetContent(content)
	m.SetSender(sender)
	m.SetRecipient(recipient)
	m.SetSubject(subject)
	return m
}

func buildMailContent(results []monitoring_check.CheckResult) string {
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

func failures(results []monitoring_check.CheckResult) []monitoring_check.CheckResult {
	failures := make([]monitoring_check.CheckResult, 0)
	for _, result := range results {
		if !result.Success() {
			failures = append(failures, result)
		}
	}
	return failures
}
