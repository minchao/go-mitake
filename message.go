package mitake

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// StatusCode of Mitake API.
type StatusCode string

func (c StatusCode) String() string {
	return statusCodeMap[c]
}

// List of Mitake API status codes.
const (
	StatusServiceError                  = StatusCode("*")
	StatusSMSTemporarilyUnavailable     = StatusCode("a")
	StatusSMSTemporarilyUnavailableB    = StatusCode("b")
	StatusUsernameRequired              = StatusCode("c")
	StatusPasswordRequired              = StatusCode("d")
	StatusUsernameOrPasswordError       = StatusCode("e")
	StatusAccountExpired                = StatusCode("f")
	StatusAccountDisabled               = StatusCode("h")
	StatusInvalidConnectionAddress      = StatusCode("k")
	StatusReachConnectUserLimit         = StatusCode("l")
	StatusChangePasswordRequired        = StatusCode("m")
	StatusPasswordExpired               = StatusCode("n")
	StatusPermissionDenied              = StatusCode("p")
	StatusServiceTemporarilyUnavailable = StatusCode("r")
	StatusAccountingFailure             = StatusCode("s")
	StatusSMSExpired                    = StatusCode("t")
	StatusSMSBodyEmpty                  = StatusCode("u")
	StatusInvalidPhoneNumber            = StatusCode("v")
	StatusQueryRecordExceededLimit      = StatusCode("w")
	StatusSMSSizeTooLarge               = StatusCode("x")
	StatusParameterError                = StatusCode("y")
	StatusNoRecord                      = StatusCode("z")
	StatusReservationForDelivery        = StatusCode("0")
	StatusCarrierAccepted               = StatusCode("1")
	StatusCarrierAccepted2              = StatusCode("2")
	StatusCarrierAccepted3              = StatusCode("3")
	StatusDelivered                     = StatusCode("4")
	StatusContentError                  = StatusCode("5")
	StatusPhoneNumberError              = StatusCode("6")
	StatusSMSDisable                    = StatusCode("7")
	StatusDeliveryTimeout               = StatusCode("8")
	StatusReservationCanceled           = StatusCode("9")
)

var statusCodeMap = map[StatusCode]string{
	StatusServiceError:                  "系統發生錯誤，請聯絡三竹資訊窗口人員",
	StatusSMSTemporarilyUnavailable:     "簡訊發送功能暫時停止服務，請稍候再試",
	StatusSMSTemporarilyUnavailableB:    "簡訊發送功能暫時停止服務，請稍候再試",
	StatusUsernameRequired:              "請輸入帳號",
	StatusPasswordRequired:              "請輸入密碼",
	StatusUsernameOrPasswordError:       "帳號、密碼錯誤",
	StatusAccountExpired:                "帳號已過期",
	StatusAccountDisabled:               "帳號已被停用",
	StatusInvalidConnectionAddress:      "無效的連線位址",
	StatusReachConnectUserLimit:         "帳號已達到同時連線數上限",
	StatusChangePasswordRequired:        "必須變更密碼，在變更密碼前，無法使用簡訊發送服務",
	StatusPasswordExpired:               "密碼已逾期，在變更密碼前，將無法使用簡訊發送服務",
	StatusPermissionDenied:              "沒有權限使用外部Http程式",
	StatusServiceTemporarilyUnavailable: "系統暫停服務，請稍後再試",
	StatusAccountingFailure:             "帳務處理失敗，無法發送簡訊",
	StatusSMSExpired:                    "簡訊已過期",
	StatusSMSBodyEmpty:                  "簡訊內容不得為空白",
	StatusInvalidPhoneNumber:            "無效的手機號碼",
	StatusQueryRecordExceededLimit:      "查詢筆數超過上限",
	StatusSMSSizeTooLarge:               "發送檔案過大，無法發送簡訊",
	StatusParameterError:                "參數錯誤",
	StatusNoRecord:                      "查無資料",
	StatusReservationForDelivery:        "預約傳送中",
	StatusCarrierAccepted:               "已送達業者",
	StatusCarrierAccepted2:              "已送達業者",
	StatusCarrierAccepted3:              "已送達業者",
	StatusDelivered:                     "已送達手機",
	StatusContentError:                  "內容有錯誤",
	StatusPhoneNumberError:              "門號有錯誤",
	StatusSMSDisable:                    "簡訊已停用",
	StatusDeliveryTimeout:               "逾時無送達",
	StatusReservationCanceled:           "預約已取消",
}

// Message represents an SMS object.
type Message struct {
	Dstaddr  string `json:"dstaddr"`  // Destination phone number
	Destname string `json:"destname"` // Destination receiver name
	Dlvtime  string `json:"dlvtime"`  // Optional, Delivery time
	Vldtime  string `json:"vldtime"`  // Optional
	Smbody   string `json:"smbody"`   // The text of the message you want to send
	Response string `json:"response"` // Optional, Callback URL to receive the delivery receipt of the message
	ClientID string `json:"clientid"` // Optional (required when bulk send), an unique identifier from client to identify SMS message
	// ObjectID string `json:"objectID"` // Optional
}

// ToINI returns the INI format string from the message fields.
func (m Message) ToINI() string {
	smbody := strings.Replace(m.Smbody, "\n", string(byte(6)), -1)

	var ini string
	ini += "dstaddr=" + m.Dstaddr + "&"
	ini += "smbody=" + smbody + "&"
	if m.Destname != "" {
		ini += "destname=" + m.Destname + "&"
	}
	if m.Dlvtime != "" {
		ini += "dlvtime=" + m.Dlvtime + "&"
	}
	if m.Vldtime != "" {
		ini += "vldtime=" + m.Vldtime + "&"
	}
	if m.Response != "" {
		ini += "response=" + m.Response + "&"
	}
	if m.ClientID != "" {
		ini += "clientid=" + m.ClientID + "&"
	}
	return ini
}

// ToBatchMessage returns the format string for multiple SMS.
func (m Message) ToBatchMessage() string {
	smbody := strings.Replace(m.Smbody, "\n", string(byte(6)), -1)

	var ini string
	ini += m.ClientID + "$$" // The document says this field is REQUIRED
	ini += m.Dstaddr + "$$"
	if m.Dlvtime != "" {
		ini += m.Dlvtime
	}
	ini += "$$"
	if m.Vldtime != "" {
		ini += m.Vldtime
	}
	ini += "$$"
	if m.Destname != "" {
		ini += m.Destname
	}
	ini += "$$"
	if m.Response != "" {
		ini += m.Response
	}
	ini += "$$"
	ini += smbody + "\n"
	return ini
}

// MessageResult represents result of send SMS.
type MessageResult struct {
	ClientID     string     `json:"clientid"`
	Msgid        string     `json:"msgid"`
	Statuscode   string     `json:"statuscode"`
	Statusstring StatusCode `json:"statusstring"`
	Duplicate    string     `json:"Duplicate"`
}

// MessageResponse represents response of send SMS.
type MessageResponse struct {
	Results      []*MessageResult
	AccountPoint int
	INI          string `json:"-"`
}

func parseMessageResponseByPattern(pattern string, body io.Reader) (*MessageResponse, error) {
	var (
		scanner  = bufio.NewScanner(transform.NewReader(body, traditionalchinese.Big5.NewDecoder()))
		response = new(MessageResponse)
		result   *MessageResult
	)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		response.INI += text + "\n"

		if matched, _ := regexp.MatchString(pattern, text); matched {
			result = new(MessageResult)
			response.Results = append(response.Results, result)
			// ClientID
			re := regexp.MustCompile(pattern)
			result.ClientID = re.FindStringSubmatch(text)[1]
		} else {
			strs := strings.Split(text, "=")
			switch strs[0] {
			case "msgid":
				result.Msgid = strs[1]
			case "statuscode":
				result.Statuscode = strs[1]
				result.Statusstring = StatusCode(strs[1])
			case "Duplicate":
				result.Duplicate = strs[1]
			case "AccountPoint":
				response.AccountPoint, _ = strconv.Atoi(strs[1])
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return response, nil
}

func parseMessageResponse(body io.Reader) (*MessageResponse, error) {
	return parseMessageResponseByPattern(`^\[(.+?)\]$`, body)
}

// MessageStatus represents status of message.
type MessageStatus struct {
	MessageResult
	StatusTime string `json:"statustime"`
}

// MessageStatusResponse represents response of query message status.
type MessageStatusResponse struct {
	Statuses []*MessageStatus
	INI      string `json:"-"`
}

func parseMessageStatusResponse(body io.Reader) (*MessageStatusResponse, error) {
	var (
		scanner  = bufio.NewScanner(transform.NewReader(body, traditionalchinese.Big5.NewDecoder()))
		response = new(MessageStatusResponse)
	)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		response.INI += text + "\n"

		strs := strings.Split(text, "\t")
		response.Statuses = append(response.Statuses, &MessageStatus{
			MessageResult: MessageResult{
				Msgid:        strs[0],
				Statuscode:   strs[1],
				Statusstring: StatusCode(strs[1]),
			},
			StatusTime: strs[2],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return response, nil
}

func parseCancelMessageStatusResponse(body io.Reader) (*MessageStatusResponse, error) {
	var (
		scanner  = bufio.NewScanner(transform.NewReader(body, traditionalchinese.Big5.NewDecoder()))
		response = new(MessageStatusResponse)
	)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		response.INI += text + "\n"

		strs := strings.Split(text, "=")
		response.Statuses = append(response.Statuses, &MessageStatus{
			MessageResult: MessageResult{
				Msgid:        strs[0],
				Statuscode:   strs[1],
				Statusstring: StatusCode(strs[1]),
			},
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return response, nil
}

// MessageReceipt represents a message delivery receipt.
type MessageReceipt struct {
	Msgid        string     `json:"msgid"`
	Dstaddr      string     `json:"dstaddr"`
	Dlvtime      string     `json:"dlvtime"`
	Donetime     string     `json:"donetime"`
	Statuscode   string     `json:"statuscode"`
	Statusstring StatusCode `json:"statusstring"`
	Statusstr    string     `json:"statusstr"`
	StatusFlag   string     `json:"StatusFlag"`
}

// ParseMessageReceipt parse an incoming Mitake callback request and return the MessageReceipt.
//
// Example usage:
//
// 	func Callback(w http.ResponseWriter, r *http.Request) {
// 		receipt, err := mitake.ParseMessageReceipt(r)
// 		if err != nil { ... }
//		// Process message receipt
// 	}
//
func ParseMessageReceipt(r *http.Request) (*MessageReceipt, error) {
	values := r.URL.Query()
	if values.Get("msgid") == "" {
		return nil, errors.New("receipt not found")
	}
	return &MessageReceipt{
		Msgid:        values.Get("msgid"),
		Dstaddr:      values.Get("dstaddr"),
		Dlvtime:      values.Get("dlvtime"),
		Donetime:     values.Get("donetime"),
		Statuscode:   values.Get("statuscode"),
		Statusstring: StatusCode(values.Get("statuscode")),
		Statusstr:    values.Get("statusstr"),
		StatusFlag:   values.Get("StatusFlag"),
	}, nil
}
