package mitake

import (
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) SendBatch(messages []Message) (*MessageResponse, error) {
	q := c.buildDefaultQuey()
	q.Set("encoding", "UTF8")
	url, _ := url.Parse("SmSendPost.asp")
	url.RawQuery = q.Encode()

	var ini string
	for i, message := range messages {
		ini += "[" + strconv.Itoa(i) + "]\n"
		ini += message.ToINI()
	}
	ini = strings.TrimSpace(ini)

	resp, err := c.Post(url.String(), "text/plain", strings.NewReader(ini))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseMessageResponse(resp.Body)
}

func (c *Client) Send(message Message) (*MessageResponse, error) {
	return c.SendBatch([]Message{message})
}

func (c *Client) QueryAccountPoint() (int, error) {
	url, _ := url.Parse("SmQueryGet.asp")
	url.RawQuery = c.buildDefaultQuey().Encode()

	resp, err := c.Get(url.String())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.Split(string(data), "=")[1])
}

func (c *Client) QueryMessageStatus(messageIds []string) (*MessageStatusResponse, error) {
	q := c.buildDefaultQuey()
	q.Set("msgid", strings.Join(messageIds, ","))

	url, _ := url.Parse("SmQueryGet.asp")
	url.RawQuery = q.Encode()

	resp, err := c.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseMessageStatusResponse(resp.Body)
}
