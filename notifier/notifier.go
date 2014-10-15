package notifier

import (
	"bytes"
	"fmt"

	"github.com/bborbe/monitoring/check"
)

func Notify(results []check.CheckResult) error {
	mailContent := buildMailContent(results)
	return sendMail(mailContent)
}

func sendMail(content string) error {
	return nil
}

func buildMailContent(results []check.CheckResult) string {
	failures := failures(results)
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Checks executed: %d\n", len(results)))
	buffer.WriteString(fmt.Sprintf("Checks failed: %d\n", len(failures)))
	for _, result := range failures {
		buffer.WriteString(fmt.Sprintf("%s - %v\n", result.Message(), result.Error()))
	}
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
