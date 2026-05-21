package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Client struct {
	conn *nats.Conn
}

func New(url string) (*Client, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) PublishJSON(subject string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return c.conn.Publish(subject, data)
}

func (c *Client) Subscribe(subject string, handler func(subject string, data []byte)) (*nats.Subscription, error) {
	return c.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Subject, msg.Data)
	})
}

func (c *Client) Close() {
	c.conn.Drain()
	c.conn.Close()
}
