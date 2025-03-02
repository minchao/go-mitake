package mitake

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type MessageParams struct {
	Encoding           string // The encoding of the message body
	ObjectID           string // Name fo the batch
	HideDeductedPoints bool   // Set to true to hide the points deducted per SMS in the response
	Message
}

func (p MessageParams) Validate() error {
	if p.Dstaddr == "" {
		return &ParameterError{Reason: "empty Dstaddr"}
	}
	if p.Smbody == "" {
		return &ParameterError{Reason: "empty Smbody"}
	}
	return nil
}

// ToData converts the message to url.Values for sending.
func (p MessageParams) ToData() url.Values {
	data := url.Values{}
	data.Set("dstaddr", p.Dstaddr)
	if p.Destname != "" {
		data.Set("destname", p.Destname)
	}
	if p.Dlvtime != "" {
		data.Set("dlvtime", p.Dlvtime)
	}
	if p.Vldtime != "" {
		data.Set("vldtime", p.Vldtime)
	}
	data.Set("smbody", p.Smbody)
	if p.Response != "" {
		data.Set("response", p.Response)
	}
	if p.ClientID != "" {
		data.Set("clientid", p.ClientID)
	}
	if p.ObjectID != "" {
		data.Set("objectID", p.ObjectID)
	}
	if !p.HideDeductedPoints {
		data.Set("smsPointFlag", "1")
	}
	return data
}

// Send sends a SMS.
func (c *Client) Send(ctx context.Context, params MessageParams) (*MessageResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	u, _ := url.Parse("b2c/mtk/SmSend")
	u.RawQuery = c.buildSendQuery(params).Encode()
	data := c.buildSendFormData(params)

	resp, err := c.Post(ctx, u.String(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseMessageResponse(resp.Body)
}

func (c *Client) buildSendQuery(params MessageParams) url.Values {
	encoding := defaultEncoding
	if params.Encoding != "" {
		encoding = params.Encoding
	}
	q := url.Values{}
	q.Set("CharsetURL", encoding)
	return q
}

func (c *Client) buildSendFormData(params MessageParams) url.Values {
	data := params.ToData()
	data.Set("username", c.username)
	data.Set("password", c.password)
	return data
}

type BatchMessagesParams struct {
	Encoding           string `json:"Encoding_postIn"` // The encoding of the message body
	ObjectID           string `json:"objectID"`        // Name fo the batch
	HideDeductedPoints bool   // Set to true to hide the points deducted per SMS in the response
	Messages           []Message
}

func (p BatchMessagesParams) Validate() error {
	if len(p.Messages) == 0 {
		return &ParameterError{Reason: "empty messages"}
	}
	for i, message := range p.Messages {
		if message.ClientID == "" {
			return &ParameterError{Reason: fmt.Sprintf("%d: empty ClientID", i)}
		}
		if message.Dstaddr == "" {
			return &ParameterError{Reason: fmt.Sprintf("%d: [%s] empty Dstaddr", i, message.ClientID)}
		}
		if message.Smbody == "" {
			return &ParameterError{Reason: fmt.Sprintf("%d: [%s] empty Smbody", i, message.ClientID)}
		}
	}
	return nil
}

func (p BatchMessagesParams) ToData() string {
	var data string
	for _, message := range p.Messages {
		data += fmt.Sprintf("%s$$%s$$%s$$%s$$%s$$%s$$%s\r\n",
			message.ClientID,
			message.Dstaddr,
			message.Dlvtime,
			message.Vldtime,
			message.Destname,
			message.Response,
			message.Smbody,
		)
	}
	return data
}

// SendBatch sends multiple SMS.
func (c *Client) SendBatch(ctx context.Context, opts BatchMessagesParams) (*MessageResponse, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	u, _ := url.Parse("b2c/mtk/SmBulkSend")
	u.RawQuery = c.buildSendBatchQuery(opts).Encode()
	data := opts.ToData()

	resp, err := c.Post(ctx, u.String(), "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseMessageResponse(resp.Body)
}

func (c *Client) buildSendBatchQuery(opts BatchMessagesParams) url.Values {
	encoding := defaultEncoding
	if opts.Encoding != "" {
		encoding = opts.Encoding
	}

	q := url.Values{}
	q.Set("username", c.username)
	q.Set("password", c.password)
	q.Set("Encoding_PostIn", encoding)
	if opts.ObjectID != "" {
		q.Set("objectID", opts.ObjectID)
	}
	if !opts.HideDeductedPoints {
		q.Set("smsPointFlag", "1")
	}
	return q
}

// MessageResult represents result of send SMS.
type MessageResult struct {
	Msgid      string
	StatusCode StatusCode
	SmsPoint   *int // Points deducted per SMS, only available when SmsPointFlag is set
}

// MessageResponse represents response of send SMS.
type MessageResponse struct {
	Results      []*MessageResult
	AccountPoint int     // The available balance after this send
	Duplicate    *string // `Y` if the message is duplicated
}

func parseMessageResponse(body io.Reader) (*MessageResponse, error) {
	var (
		scanner  = bufio.NewScanner(body)
		re       = regexp.MustCompile(`^\[(.+?)]$`)
		response = new(MessageResponse)
		result   *MessageResult
	)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if matched := re.MatchString(text); matched {
			result = new(MessageResult)
			response.Results = append(response.Results, result)
		} else {
			if result == nil {
				return nil, &UnexpectedResponseError{Reason: "no clientid"}
			}
			s := strings.Split(text, "=")
			if len(s) != 2 {
				return nil, &UnexpectedResponseError{Reason: "invalid key value pair"}
			}

			switch s[0] {
			case "msgid":
				result.Msgid = s[1]
			case "statuscode":
				result.StatusCode = StatusCode(s[1])
			case "smsPoint":
				point, _ := strconv.Atoi(s[1])
				result.SmsPoint = Ptr(point)
			case "AccountPoint":
				response.AccountPoint, _ = strconv.Atoi(s[1])
			case "Duplicate":
				response.Duplicate = Ptr(s[1])
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if result == nil {
		return nil, &UnexpectedResponseError{Reason: "invalid response"}
	}
	return response, nil
}

// MessageStatusParams represents parameters of query message status.
type MessageStatusParams struct {
	MessageIDs         []string
	HideDeductedPoints bool
}

// MessageStatus represents status of message.
type MessageStatus struct {
	MessageResult
	StatusTime string
}

// MessageStatusResponse represents response of query message status.
type MessageStatusResponse struct {
	Statuses []*MessageStatus
}

// QueryMessageStatus fetch the status of specific messages.
func (c *Client) QueryMessageStatus(ctx context.Context, params MessageStatusParams) (*MessageStatusResponse, error) {
	q := c.buildDefaultQuery()
	q.Set("msgid", strings.Join(params.MessageIDs, ","))
	if !params.HideDeductedPoints {
		q.Set("smsPointFlag", "1")
	}

	u, _ := url.Parse("/b2c/mtk/SmQuery")
	u.RawQuery = q.Encode()

	resp, err := c.Get(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseMessageStatusResponse(resp.Body)
}

func parseMessageStatusResponse(body io.Reader) (*MessageStatusResponse, error) {
	var (
		scanner  = bufio.NewScanner(body)
		response = new(MessageStatusResponse)
	)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		s := strings.Split(text, "\t")
		messageStatus := &MessageStatus{
			MessageResult: MessageResult{
				Msgid:      s[0],
				StatusCode: StatusCode(s[1]),
			},
			StatusTime: s[2],
		}
		if len(s) == 4 {
			point, _ := strconv.Atoi(s[3])
			messageStatus.SmsPoint = Ptr(point)
		}
		response.Statuses = append(response.Statuses, messageStatus)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return response, nil
}

// QueryAccountPoint retrieves your account balance.
func (c *Client) QueryAccountPoint(ctx context.Context) (int, error) {
	u, _ := url.Parse("b2c/mtk/SmQuery")
	u.RawQuery = c.buildDefaultQuery().Encode()

	resp, err := c.Get(ctx, u.String())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.Split(string(data), "=")[1])
}

// CanceledMessage represents the canceled message.
type CanceledMessage struct {
	Msgid      string
	StatusCode StatusCode
}

// CancelScheduledMessages cancels scheduled messages.
func (c *Client) CancelScheduledMessages(ctx context.Context, messageIDs []string) ([]*CanceledMessage, error) {
	q := c.buildDefaultQuery()
	q.Set("msgid", strings.Join(messageIDs, ","))

	u, _ := url.Parse("b2c/mtk/SmCancel")
	u.RawQuery = q.Encode()

	resp, err := c.Get(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseCancelScheduledMessagesResponse(resp.Body)
}

func parseCancelScheduledMessagesResponse(body io.Reader) ([]*CanceledMessage, error) {
	var (
		scanner  = bufio.NewScanner(body)
		messages = make([]*CanceledMessage, 0)
	)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		s := strings.Split(text, "=")
		messages = append(messages, &CanceledMessage{
			Msgid:      s[0],
			StatusCode: StatusCode(s[1]),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
