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

	rel, err := url.Parse("SmSendPost.asp")
	if err != nil {
		return nil, err
	}
	url := c.BaseURL.ResolveReference(rel)
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

// SendBatchLong sends long SMS.
func (c *Client) SendLongMessageBatch(messages []Message) (*MessageResponse, error) {
	q := c.buildDefaultQuery()
	q.Set("Encoding_PostIn", "UTF8")

	rel, err := url.Parse("SpLmPost")
	if err != nil {
		return nil, err
	}
	url := c.LongMessageURL.ResolveReference(rel)
	url.RawQuery = q.Encode()

	var ini string
	for _, message := range messages {
		ini += message.ID + "$$"
		ini += message.ToLM()
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
	rel, err := url.Parse("SmQueryGet.asp")
	if err != nil {
		return 0, err
	}
	url := c.BaseURL.ResolveReference(rel)
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

	rel, err := url.Parse("SmQueryGet.asp")
	if err != nil {
		return nil, err
	}
	url := c.BaseURL.ResolveReference(rel)
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

	rel, err := url.Parse("SmCancel.asp")
	if err != nil {
		return nil, err
	}
	url := c.BaseURL.ResolveReference(rel)
	url.RawQuery = q.Encode()

	resp, err := c.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseCancelMessageStatusResponse(resp.Body)
}
