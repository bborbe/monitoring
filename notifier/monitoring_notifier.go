package notifier

import (
	"bytes"
	"fmt"

	"github.com/bborbe/log"
	"github.com/bborbe/mail"
	"github.com/bborbe/mail/message"
	"github.com/bborbe/monitoring/check"
)

var logger = log.DefaultLogger

type notifier struct {
	mailer    mail.Mailer
	sender    string
	recipient string
}

type Notifier interface {
	Notify(results []check.CheckResult) error
}

func New(mailer mail.Mailer, sender string, recipient string) *notifier {
	n := new(notifier)
	n.mailer = mailer
	n.sender = sender
	n.recipient = recipient
	return n
}

func (n *notifier) Notify(results []check.CheckResult) error {
	logger.Debug("notify results")
	mailContent := buildMailContent(results)
	message := buildMessage(n.sender, n.recipient, mailContent)
	err := n.mailer.Send(message)
	logger.Debug("mail sent")
	return err
}

func buildMessage(sender string, recipient string, content string) mail.Message {
	m := message.New()
	m.SetContent(content)
	m.SetSender(sender)
	m.SetRecipient(recipient)
	return m
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
