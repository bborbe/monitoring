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
	"github.com/bborbe/monitoring/runner"
)

var logger = log.DefaultLogger

func main() {
	defer logger.Close()
	logLevelPtr := flag.Int("loglevel", log.OFF, "int")
	flag.Parse()
	logger.SetLevelThreshold(*logLevelPtr)
	logger.Debugf("set log level to %s", log.LogLevelToString(*logLevelPtr))

	writer := os.Stdout
	err := do(writer)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(writer io.Writer) error {
	var err error
	fmt.Fprintf(writer, "check started\n")
	c := configuration.New()
	resultChannel := runner.Run(c.Checks())
	results := make([]check.CheckResult, 0)
	hasError := false
	for result := range resultChannel {
		if result.Success() {
			fmt.Fprintf(writer, "[OK]   %s\n", result.Message())
		} else {
			fmt.Fprintf(writer, "[FAIL] %s - %v\n", result.Message(), result.Error())
			hasError = true
		}
		results = append(results, result)
	}
	if hasError {
		err = notifier.Notify(results)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(writer, "check finished\n")
	return err
}
