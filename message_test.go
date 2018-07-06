package mitake

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestMessage_ToINI(t *testing.T) {
	message1 := Message{
		Dstaddr:  "0987654321",
		Smbody:   "Test",
		Dlvtime:  "20170101010000",
		Vldtime:  "20170101012300",
		Response: "https://example.com/callback",
	}
	want1 := "dstaddr=0987654321\nsmbody=Test\ndlvtime=20170101010000\nvldtime=20170101012300\nresponse=https://example.com/callback\n"
	if got := message1.ToINI(); got != want1 {
		t.Errorf("Message INI is %v, want %v", got, want1)
	}
	message2 := Message{
		Dstaddr: "0987654321",
		Smbody:  "Test",
	}
	want2 := "dstaddr=0987654321\nsmbody=Test\n"
	if got := message2.ToINI(); got != want2 {
		t.Errorf("Message INI is %v, want %v", got, want1)
	}
}

func TestMessage_ToLM(t *testing.T) {
	message1 := Message{
		Dstaddr:  "0987654321",
		Destname: "Bob",
		Smbody:   "Test",
		Dlvtime:  "20170101010000",
		Vldtime:  "20170101012300",
		Response: "https://example.com/callback",
	}
	want1 := "0987654321$$20170101010000$$20170101012300$$Bob$$https://example.com/callback$$Test\n"
	if got := message1.ToLM(); got != want1 {
		t.Errorf("Message LM is %v, want %v", got, want1)
	}
	message2 := Message{
		Dstaddr: "0987654321",
		Smbody:  "Test",
	}
	want2 := "0987654321$$$$$$$$$$Test\n"
	if got := message2.ToLM(); got != want2 {
		t.Errorf("Message LM is %v, want %v", got, want1)
	}
}

func Test_parseMessageResponse(t *testing.T) {
	body := strings.NewReader(`[0]
msgid=1010079522
statuscode=1
[1]
msgid=1010079523
statuscode=4
AccountPoint=98`)
	resp, err := parseMessageResponse(body)
	if err != nil {
		t.Errorf("parseMessageResponse returned unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Errorf("MessageResponse.Result len is %d, want %d", len(resp.Results), 2)
	}
	if resp.AccountPoint != 98 {
		t.Errorf("MessageResponse.AccountPoint is %d, want %d", resp.AccountPoint, 98)
	}

	want := []*MessageResult{
		{
			Msgid:        "1010079522",
			Statuscode:   "1",
			Statusstring: StatusCode("1"),
		},
		{
			Msgid:        "1010079523",
			Statuscode:   "4",
			Statusstring: StatusCode("4"),
		},
	}
	if !reflect.DeepEqual(resp.Results, want) {
		t.Errorf("MessageResult returned %+v, want %+v", resp.Results, want)
	}
}

func Test_parseLongMessageResponse(t *testing.T) {
	body := strings.NewReader(`[2ks8k828j5]
msgid=#1010079522
statuscode=1
[19ke8ks83]
msgid=#1010079523
statuscode=4
AccountPoint=98`)
	resp, err := parseLongMessageResponse(body)
	if err != nil {
		t.Errorf("parseMessageResponse returned unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Errorf("MessageResponse.Result len is %d, want %d", len(resp.Results), 2)
	}
	if resp.AccountPoint != 98 {
		t.Errorf("MessageResponse.AccountPoint is %d, want %d", resp.AccountPoint, 98)
	}

	want := []*MessageResult{
		{
			Msgid:        "#1010079522",
			Statuscode:   "1",
			Statusstring: StatusCode("1"),
		},
		{
			Msgid:        "#1010079523",
			Statuscode:   "4",
			Statusstring: StatusCode("4"),
		},
	}
	if !reflect.DeepEqual(resp.Results, want) {
		t.Errorf("MessageResult returned %+v, want %+v", resp.Results, want)
	}
}

func Test_parseMessageStatusResponse(t *testing.T) {
	body := strings.NewReader(`1010079522	1	20170101010010
1010079523	4	20170101010011`)
	resp, err := parseMessageStatusResponse(body)
	if err != nil {
		t.Errorf("parseMessageStatusResponse returned unexpected error: %v", err)
	}
	if len(resp.Statuses) != 2 {
		t.Errorf("MessageStatusResponse.Statuses len is %d, want %d", len(resp.Statuses), 2)
	}

	want := []*MessageStatus{
		{
			MessageResult: MessageResult{
				Msgid:        "1010079522",
				Statuscode:   "1",
				Statusstring: StatusCode("1"),
			},
			StatusTime: "20170101010010",
		},
		{
			MessageResult: MessageResult{
				Msgid:        "1010079523",
				Statuscode:   "4",
				Statusstring: StatusCode("4"),
			},
			StatusTime: "20170101010011",
		},
	}
	if !reflect.DeepEqual(resp.Statuses, want) {
		t.Errorf("MessageStatus returned %+v, want %+v", resp.Statuses, want)
	}
}

func Test_parseMessageCancelStatusResponse(t *testing.T) {
	body := strings.NewReader(`1010079522=8
1010079523=9`)
	resp, err := parseCancelMessageStatusResponse(body)
	if err != nil {
		t.Errorf("parseMessageStatusResponse returned unexpected error: %v", err)
	}
	if len(resp.Statuses) != 2 {
		t.Errorf("MessageStatusResponse.Statuses len is %d, want %d", len(resp.Statuses), 2)
	}

	want := []*MessageStatus{
		{
			MessageResult: MessageResult{
				Msgid:        "1010079522",
				Statuscode:   "8",
				Statusstring: StatusCode("8"),
			},
		},
		{
			MessageResult: MessageResult{
				Msgid:        "1010079523",
				Statuscode:   "9",
				Statusstring: StatusCode("9"),
			},
		},
	}
	if !reflect.DeepEqual(resp.Statuses, want) {
		t.Errorf("MessageStatus returned %+v, want %+v", resp.Statuses, want)
	}
}

func Test_ParseMessageReceipt(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		receipt, err := ParseMessageReceipt(r)
		if err != nil {
			t.Errorf("ParseMessageReceipt returned unexpected error: %v", err)
			return
		}

		want := &MessageReceipt{
			Msgid:        "8091234567",
			Dstaddr:      "09001234567",
			Dlvtime:      "20060810125612",
			Donetime:     "20060810165612",
			Statuscode:   "0",
			Statusstring: StatusCode("0"),
			Statusstr:    "DELIVRD",
			StatusFlag:   "4",
		}

		if !reflect.DeepEqual(receipt, want) {
			t.Errorf("Message received: %v, want %v", receipt, want)
		}
	})

	// Simulate the mitake server response.
	client.Get("/callback" +
		"?msgid=8091234567" +
		"&dstaddr=09001234567" +
		"&dlvtime=20060810125612" +
		"&donetime=20060810165612" +
		"&statusstr=DELIVRD" +
		"&statuscode=0" +
		"&StatusFlag=4")
}
