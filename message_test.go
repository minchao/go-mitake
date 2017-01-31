package mitake

import (
	"reflect"
	"strings"
	"testing"
)

func TestMessage_ToINI(t *testing.T) {
	message1 := Message{
		Dstaddr: "0987654321",
		Smbody:  "Test",
		Dlvtime: "20170101010000",
		Vldtime: "20170101012300",
	}
	want1 := "dstaddr=0987654321\nsmbody=Test\ndlvtime=20170101010000\nvldtime=20170101012300\n"
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
			Msgid:      "1010079522",
			Statuscode: StatusCode("1"),
		},
		{
			Msgid:      "1010079523",
			Statuscode: StatusCode("4"),
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
				Msgid:      "1010079522",
				Statuscode: StatusCode("1"),
			},
			StatusTime: "20170101010010",
		},
		{
			MessageResult: MessageResult{
				Msgid:      "1010079523",
				Statuscode: StatusCode("4"),
			},
			StatusTime: "20170101010011",
		},
	}
	if !reflect.DeepEqual(resp.Statuses, want) {
		t.Errorf("MessageStatus returned %+v, want %+v", resp.Statuses, want)
	}
}
