package mock

import "github.com/bborbe/mailer"

type mailerMock struct {
	Error   error
	Message mailer.Message
	Counter int
}

func New() *mailerMock {
	m := new(mailerMock)
	m.Counter = 0
	return m
}

func (m *mailerMock) Send(message mailer.Message) error {
	m.Message = message
	m.Counter++
	return m.Error
}
