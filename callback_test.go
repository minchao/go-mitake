package mitake

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func Test_ParseMessageReceipt(t *testing.T) {
	client, mux, teardown := setup()
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
			t.Errorf("MessageTemp received: %v, want %v", receipt, want)
		}
	})

	// Simulate the mitake server response.
	_, _ = client.Get(context.Background(),
		"/callback"+
			"?msgid=8091234567"+
			"&dstaddr=09001234567"+
			"&dlvtime=20060810125612"+
			"&donetime=20060810165612"+
			"&statusstr=DELIVRD"+
			"&statuscode=0"+
			"&StatusFlag=4",
	)
}

func Test_ParseMessageReceipt_emptyMsgid(t *testing.T) {
	_, err := ParseMessageReceipt(&http.Request{URL: &url.URL{}})
	if err == nil {
		t.Fatal("ParseMessageReceipt did not return an error")
	}
	if err.Error() != "receipt not found" {
		t.Errorf("ParseMessageReceipt returned unexpected error: %v", err)
	}
}
