package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"runtime"
	"sync"

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

	var wg sync.WaitGroup

	maxConcurrency := runtime.NumCPU() * 2
	throttle := make(chan bool, maxConcurrency)

	c := configuration.New()
	for _, check := range c.Checks() {
		c := check
		wg.Add(1)
		go func() {
			throttle <- true
			result := c.Check()
			<-throttle
			if result.Success() {
				fmt.Fprintf(writer, "[OK]   %s\n", result.Message())
			} else {
				fmt.Fprintf(writer, "[FAIL] %s - %v\n", result.Message(), result.Error())
			}
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Fprintf(writer, "check finished\n")
	return nil
}
