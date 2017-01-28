package mitake

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type StatusCode string

func (c StatusCode) String() string {
	return statusCodeMap[c]
}

const (
	StatusSTAR = StatusCode("*")
	StatusA    = StatusCode("a")
	StatusB    = StatusCode("b")
	StatusC    = StatusCode("c")
	StatusD    = StatusCode("d")
	StatusE    = StatusCode("e")
	StatusF    = StatusCode("f")
	StatusH    = StatusCode("g")
	StatusK    = StatusCode("k")
	StatusM    = StatusCode("m")
	StatusN    = StatusCode("n")
	StatusP    = StatusCode("p")
	StatusR    = StatusCode("r")
	StatusS    = StatusCode("s")
	StatusT    = StatusCode("t")
	StatusU    = StatusCode("u")
	StatusV    = StatusCode("v")
	Status0    = StatusCode("0")
	Status1    = StatusCode("1")
	Status2    = StatusCode("2")
	Status3    = StatusCode("3")
	Status4    = StatusCode("4")
	Status5    = StatusCode("5")
	Status6    = StatusCode("6")
	Status7    = StatusCode("7")
	Status8    = StatusCode("8")
	Status9    = StatusCode("9")
)

var statusCodeMap = map[StatusCode]string{
	StatusSTAR: "系統發生錯誤，請聯絡三竹資訊窗口人員",
	StatusA:    "簡訊發送功能暫時停止服務，請稍候再試",
	StatusB:    "簡訊發送功能暫時停止服務，請稍候再試",
	StatusC:    "請輸入帳號",
	StatusD:    "請輸入密碼",
	StatusE:    "帳號、密碼錯誤",
	StatusF:    "帳號已過期",
	StatusH:    "帳號已被停用",
	StatusK:    "無效的連線位址",
	StatusM:    "必須變更密碼，在變更密碼前，無法使用簡訊發送服務",
	StatusN:    "密碼已逾期，在變更密碼前，將無法使用簡訊發送服務",
	StatusP:    "沒有權限使用外部Http程式",
	StatusR:    "系統暫停服務，請稍後再試",
	StatusS:    "帳務處理失敗，無法發送簡訊",
	StatusT:    "簡訊已過期",
	StatusU:    "簡訊內容不得為空白",
	StatusV:    "無效的手機號碼",
	Status0:    "預約傳送中",
	Status1:    "已送達業者",
	Status2:    "已送達業者",
	Status3:    "已送達業者",
	Status4:    "已送達手機",
	Status5:    "內容有錯誤",
	Status6:    "門號有錯誤",
	Status7:    "簡訊已停用",
	Status8:    "逾時無送達",
	Status9:    "預約已取消",
}

type Message struct {
	Dstaddr string `json:"dstaddr"` // Destination phone number
	Smbody  string `json:"smbody"`  // The text of the message you want to send
	Dlvtime string `json:"dlvtime"` // Optional, Delivery time
	Vldtime string `json:"vldtime"` // Optional
}

// ToINI returns the INI format string from the message fields.
func (m Message) ToINI() string {
	var ini string
	ini += "dstaddr=" + m.Dstaddr + "\n"
	ini += "smbody=" + m.Smbody + "\n"
	if m.Dlvtime != "" {
		ini += "dlvtime=" + m.Dlvtime + "\n"
	}
	if m.Vldtime != "" {
		ini += "vldtime=" + m.Vldtime + "\n"
	}
	return ini
}

type MessageResult struct {
	Msgid      string     `json:"msgid"`
	Statuscode StatusCode `json:"statuscode"`
}

type MessageResponse struct {
	Results      []*MessageResult
	AccountPoint int
}

func parseMessageResponse(body io.Reader) (*MessageResponse, error) {
	var (
		scanner  = bufio.NewScanner(body)
		response = new(MessageResponse)
		result   *MessageResult
	)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if matched, _ := regexp.MatchString(`^\[\d+]$`, text); matched {
			result = new(MessageResult)
			response.Results = append(response.Results, result)
		} else {
			strs := strings.Split(text, "=")
			switch strs[0] {
			case "msgid":
				result.Msgid = strs[1]
			case "statuscode":
				result.Statuscode = StatusCode(strs[1])
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

type MessageStatus struct {
	MessageResult
	StatusTime string `json:"statustime"`
}

type MessageStatusResponse struct {
	Statuses []*MessageStatus
}

func parseMessageStatusResponse(body io.Reader) (*MessageStatusResponse, error) {
	var (
		scanner  = bufio.NewScanner(body)
		response = new(MessageStatusResponse)
	)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		strs := strings.Split(text, "\t")
		response.Statuses = append(response.Statuses, &MessageStatus{
			MessageResult: MessageResult{
				Msgid:      strs[0],
				Statuscode: StatusCode(strs[1]),
			},
			StatusTime: strs[2],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return response, nil
}

// Message delivery receipt.
type MessageReceipt struct {
	Msgid      string     `json:"msgid"`
	Dstaddr    string     `json:"dstaddr"`
	Dlvtime    string     `json:"dlvtime"`
	Donetime   string     `json:"donetime"`
	Statuscode StatusCode `json:"statuscode"`
	Statusstr  string     `json:"statusstr"`
	StatusFlag string     `json:"StatusFlag"`
}

// Parse an incoming Mitake callback request and return the MessageReceipt.
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
		Msgid:      values.Get("msgid"),
		Dstaddr:    values.Get("dstaddr"),
		Dlvtime:    values.Get("dlvtime"),
		Donetime:   values.Get("donetime"),
		Statuscode: StatusCode(values.Get("statuscode")),
		Statusstr:  values.Get("statusstr"),
		StatusFlag: values.Get("StatusFlag"),
	}, nil
}
