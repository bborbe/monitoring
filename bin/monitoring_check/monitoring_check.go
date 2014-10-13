package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/configuration"
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
	fmt.Fprintf(writer, "check started\n")

	c := configuration.New()
	for _, check := range c.Checks() {
		result := check.Check()
		if result.Success() {
			fmt.Fprintf(writer, "[OK]   %s\n", result.Message())
		} else {
			fmt.Fprintf(writer, "[FAIL] %s\n", result.Message())
		}
	}
	fmt.Fprintf(writer, "check finished\n")
	return nil
}
