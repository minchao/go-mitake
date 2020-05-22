package mitake

import (
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

// Send an SMS.
func (c *Client) Send(message Message) (*MessageResponse, error) {
	q := c.buildDefaultQuery()
	q.Set("CharsetURL", "UTF8")

	u, _ := url.Parse("/api/mtk/SmSend")
	u.RawQuery = q.Encode()

	var ini string
	ini = strings.TrimSpace(message.ToINI())
	resp, err := c.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(ini))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseMessageResponse(resp.Body)
}

// SendBatch sends multiple SMS.
func (c *Client) SendBatch(messages []Message) (*MessageResponse, error) {
	q := c.buildDefaultQuery()
	q.Set("Encoding_PostIn", "UTF8")

	u, _ := url.Parse("/api/mtk/SmBulkSend")
	u.RawQuery = q.Encode()

	var ini string
	for _, message := range messages {
		ini += message.ToBatchMessage()
	}
	ini = strings.TrimSpace(ini)

	resp, err := c.Post(u.String(), "text/plain", strings.NewReader(ini))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseMessageResponse(resp.Body)
}

// QueryAccountPoint retrieves your account balance.
func (c *Client) QueryAccountPoint() (int, error) {
	u, _ := url.Parse("/api/mtk/SmQuery")
	u.RawQuery = c.buildDefaultQuery().Encode()

	resp, err := c.Get(u.String())
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

	u, _ := url.Parse("/api/mtk/SmQuery")
	u.RawQuery = q.Encode()

	resp, err := c.Get(u.String())
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

	u, _ := url.Parse("/api/mtk/SmCancel")
	u.RawQuery = q.Encode()

	resp, err := c.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseCancelMessageStatusResponse(resp.Body)
}
