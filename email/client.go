package mail

import (
	"errors"

	"gopkg.in/gomail.v2"
)

type Client struct {
	dialer *gomail.Dialer
}

func New(account, passwd, smtp string, port int) *Client {
	return &Client{
		dialer: gomail.NewDialer(smtp, port, account, passwd),
	}
}

func (e *Client) Send(options ...Option) error {
	message, err := Message(options...)
	if err != nil {
		return err
	}
	return e.SendMore(message)
}
func (e *Client) SendMore(messages ...*gomail.Message) error {
	if len(messages) <= 0 {
		return errors.New("message is empty")
	}
	return e.dialer.DialAndSend(messages...)
}
