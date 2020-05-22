package mitake

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_SendBatch(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/mtk/SmBulkSend", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testINI(t, r, `[0]
dstaddr=0987654321
smbody=Test 1
[1]
dstaddr=0987654322
smbody=Test 2`)
		_, _ = fmt.Fprint(w, `[0]
msgid=1010079522
statuscode=1
[1]
msgid=1010079523
statuscode=4
AccountPoint=98`)
	})

	messages := []Message{
		{
			ClientID: "0",
			Dstaddr:  "0987654321",
			Smbody:   "Test 1",
		},
		{
			ClientID: "1",
			Dstaddr:  "0987654322",
			Smbody:   "Test 2",
		},
	}

	resp, err := client.SendBatch(messages)
	if err != nil {
		t.Errorf("SendBatch returned unexpected error: %v", err)
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
		t.Errorf("SendBatch returned %+v, want %+v", resp.Results, want)
	}
}

func TestClient_Send(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/mtk/SmSend", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testINI(t, r, `[0]
dstaddr=0987654321
smbody=Test 1`)
		_, _ = fmt.Fprint(w, `[0]
msgid=1010079522
statuscode=1
AccountPoint=99`)
	})

	resp, err := client.Send(
		Message{
			Dstaddr: "0987654321",
			Smbody:  "Test 1",
		},
	)
	if err != nil {
		t.Errorf("Send returned unexpected error: %v", err)
	}

	want := []*MessageResult{
		{
			Msgid:        "1010079522",
			Statuscode:   "1",
			Statusstring: StatusCode("1"),
		},
	}
	if !reflect.DeepEqual(resp.Results, want) {
		t.Errorf("Send returned %+v, want %+v", resp.Results, want)
	}
}

func TestClient_QueryAccountPoint(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/mtk/SmQuery", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `AccountPoint=100`)
	})

	ap, err := client.QueryAccountPoint()
	if err != nil {
		t.Errorf("QueryAccountPoint returned unexpected error: %v", err)
	}
	if ap != 100 {
		t.Errorf("QueryAccountPoint returned %+v, want %+v", ap, 100)
	}
}

func TestClient_QueryMessageStatus(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/mtk/SmQuery", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `1010079522	1	20170101010010
1010079523	4	20170101010011`)
	})

	resp, err := client.QueryMessageStatus([]string{"1010079522", "1010079523"})
	if err != nil {
		t.Errorf("QueryMessageStatus returned unexpected error: %v", err)
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
		t.Errorf("QueryMessageStatus returned %+v, want %+v", resp.Statuses, want)
	}
}

func TestClient_CancelMessage(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/mtk/SmCancel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `1010079522=8
1010079523=9`)
	})

	resp, err := client.CancelMessageStatus([]string{"1010079522", "1010079523"})
	if err != nil {
		t.Errorf("QueryMessageStatus returned unexpected error: %v", err)
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
		t.Errorf("QueryMessageStatus returned %+v, want %+v", resp.Statuses, want)
	}
}
