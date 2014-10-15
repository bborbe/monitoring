package main

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/io"
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/configuration"
)

func TestDoEmpty(t *testing.T) {
	writer := io.NewWriter()
	err := do(writer, NewConfigurationDummy(make([]check.Check, 0)), new(mailConfig))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.Content(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(writer.Content()) > 0, Is(true))
	if err != nil {
		t.Fatal(err)
	}
}

type configurationDummy struct {
	checks []check.Check
}

func NewConfigurationDummy(checks []check.Check) configuration.Configuration {
	c := new(configurationDummy)
	c.checks = checks
	return c
}

func (c *configurationDummy) Checks() []check.Check { return c.checks }
