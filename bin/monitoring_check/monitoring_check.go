package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/configuration"
	"github.com/bborbe/monitoring/runner/all"
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
	runner := all.New()
	results := runner.Run(c)
	for result := range results {
		if result.Success() {
			fmt.Fprintf(writer, "[OK]   %s\n", result.Message())
		} else {
			fmt.Fprintf(writer, "[FAIL] %s - %v\n", result.Message(), result.Error())
		}
	}
	fmt.Fprintf(writer, "check finished\n")
	return err
}
