package mitake

import (
	"errors"
	"net/http"
)

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
//	func Callback(w http.ResponseWriter, r *http.Request) {
//		receipt, err := mitake.ParseMessageReceipt(r)
//		if err != nil { ... }
//		// Process message receipt
//	}
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
