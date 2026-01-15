package mail

import (
	"errors"

	"gopkg.in/gomail.v2"
)

type Client interface {
	Send(options ...Option) error
	SendMore(messages ...*gomail.Message) error
}

type client struct {
	dialer *gomail.Dialer
}

func New(account, passwd, smtp string, port int) Client {
	return &client{
		dialer: gomail.NewDialer(smtp, port, account, passwd),
	}
}

func (e *client) Send(options ...Option) error {
	message, err := Message(options...)
	if err != nil {
		return err
	}
	return e.SendMore(message)
}
func (e *client) SendMore(messages ...*gomail.Message) error {
	if len(messages) <= 0 {
		return errors.New("message is empty")
	}
	return e.dialer.DialAndSend(messages...)
}
