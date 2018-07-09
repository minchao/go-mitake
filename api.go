package mitake

import (
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

// SendBatch sends multiple SMS.
func (c *Client) SendBatch(messages []Message) (*MessageResponse, error) {
	q := c.buildDefaultQuery()
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

// SendLongMessageBatch sends multiple long message SMS
func (c *Client) SendLongMessageBatch(messages []Message) (*MessageResponse, error) {
	q := c.buildDefaultQuery()
	q.Set("Encoding_PostIn", "UTF8")

	url := *c.LongMessageBaseURL
	url.Path = "SpLmPost"
	url.RawQuery = q.Encode()

	var ini string
	for _, message := range messages {
		ini += message.ID + "$$"
		ini += message.ToLongMessage()
	}
	ini = strings.TrimSpace(ini)

	resp, err := c.Post(url.String(), "text/plain", strings.NewReader(ini))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseLongMessageResponse(resp.Body)
}

// Send an SMS.
func (c *Client) Send(message Message) (*MessageResponse, error) {
	return c.SendBatch([]Message{message})
}

// SendLM an Long SMS.
func (c *Client) SendLongMessage(message Message) (*MessageResponse, error) {
	return c.SendLongMessageBatch([]Message{message})
}

// QueryAccountPoint retrieves your account balance.
func (c *Client) QueryAccountPoint() (int, error) {
	url, _ := url.Parse("SmQueryGet.asp")
	url.RawQuery = c.buildDefaultQuery().Encode()

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

// QueryMessageStatus fetch the status of specific messages.
func (c *Client) QueryMessageStatus(messageIds []string) (*MessageStatusResponse, error) {
	q := c.buildDefaultQuery()
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

// CancelMessageStatus cancel the specific messages.
func (c *Client) CancelMessageStatus(messageIds []string) (*MessageStatusResponse, error) {
	q := c.buildDefaultQuery()
	q.Set("msgid", strings.Join(messageIds, ","))

	url, _ := url.Parse("SmCancel.asp")
	url.RawQuery = q.Encode()

	resp, err := c.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseCancelMessageStatusResponse(resp.Body)
}
