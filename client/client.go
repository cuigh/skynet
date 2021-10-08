package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/config"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/skynet/contract"
)

type Client struct {
	skynetAddress string
	skynetToken   string
}

type Option func(c *Client)

func WithAddress(addr string) Option {
	return func(c *Client) {
		c.skynetAddress = addr
	}
}

func WithToken(token string) Option {
	return func(c *Client) {
		c.skynetToken = token
	}
}

func NewClient(opts ...Option) *Client {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	if c.skynetAddress == "" {
		c.skynetAddress = strings.TrimRight(config.GetString("skynet.address"), "/")
	}
	if c.skynetToken == "" {
		c.skynetToken = config.GetString("skynet.token")
	}
	return c
}

// Execute calls Skynet to dispatch task
func (c *Client) Execute(param contract.ExecuteParam) error {
	return c.do("/api/task/execute", param)
}

// Notify sends execute result of job to Skynet
func (c *Client) Notify(param contract.NotifyParam) error {
	return c.do("/api/task/notify", param)
}

func (c *Client) do(path string, args interface{}) error {
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.skynetAddress+path, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set(web.HeaderContentType, web.MIMEApplicationJSONCharsetUTF8)
	if c.skynetToken != "" {
		// TODO: throw error if token is missing
		req.Header.Set(web.HeaderAuthorization, "Bearer "+c.skynetToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	result := contract.Result{}
	err = d.Decode(&result)
	if err != nil {
		return err
	} else if result.Code != contract.CodeSuccess {
		return errors.Coded(result.Code, result.Info)
	}
	return nil
}

func init() {
	container.Put(func() *Client { return NewClient() }, container.Name("skynet.client"))
}
