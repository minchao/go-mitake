package mitake

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestMessageOptions_Validate(t *testing.T) {
	testCases := []struct {
		params   MessageParams
		expected error
	}{
		{
			expected: &ParameterError{Reason: "empty Dstaddr"},
		},
		{
			params: MessageParams{
				Message: Message{
					Smbody: "Hello, 世界",
				},
			},
			expected: &ParameterError{Reason: "empty Dstaddr"},
		},
		{
			params: MessageParams{
				Message: Message{
					Dstaddr: "0987654321",
				},
			},
			expected: &ParameterError{Reason: "empty Smbody"},
		},
		{
			params: MessageParams{
				Message: Message{
					Dstaddr: "0987654321",
					Smbody:  "Hello, 世界",
				},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case=%d", i), func(t *testing.T) {
			err := tc.params.Validate()

			if !errors.Is(err, tc.expected) {
				t.Errorf("Validate returned %v, want %v", err, tc.expected)
			}
		})
	}
}

func TestMessageOptions_ToData(t *testing.T) {
	testCases := []struct {
		params   MessageParams
		expected url.Values
	}{
		{
			expected: url.Values{
				"dstaddr":      []string{""},
				"smbody":       []string{""},
				"smsPointFlag": []string{"1"},
			},
		},
		{
			params: MessageParams{
				Message: Message{
					Dstaddr: "0987654321",
					Smbody:  "Hello, 世界",
				},
			},
			expected: url.Values{
				"dstaddr":      []string{"0987654321"},
				"smbody":       []string{"Hello, 世界"},
				"smsPointFlag": []string{"1"},
			},
		},
		{
			params: MessageParams{
				Message: Message{
					Dstaddr: "0987654321",
					Dlvtime: "20170101010000",
					Smbody:  "Hello, 世界",
				},
			},
			expected: url.Values{
				"dlvtime":      []string{"20170101010000"},
				"dstaddr":      []string{"0987654321"},
				"smbody":       []string{"Hello, 世界"},
				"smsPointFlag": []string{"1"},
			},
		},
		{
			params: MessageParams{
				Message: Message{
					Dstaddr: "0987654321",
					Vldtime: "20170101013000",
					Smbody:  "Hello, 世界",
				},
			},
			expected: url.Values{
				"dstaddr":      []string{"0987654321"},
				"smbody":       []string{"Hello, 世界"},
				"vldtime":      []string{"20170101013000"},
				"smsPointFlag": []string{"1"},
			},
		},
		{
			params: MessageParams{
				Message: Message{
					Dstaddr:  "0987654321",
					Destname: "Bob",
					Smbody:   "Hello, 世界",
				},
			},
			expected: url.Values{
				"destname":     []string{"Bob"},
				"dstaddr":      []string{"0987654321"},
				"smbody":       []string{"Hello, 世界"},
				"smsPointFlag": []string{"1"},
			},
		},
		{
			params: MessageParams{
				Message: Message{
					Dstaddr:  "0987654321",
					Smbody:   "Hello, 世界",
					Response: "https://example.com/callback",
				},
			},
			expected: url.Values{
				"dstaddr":      []string{"0987654321"},
				"response":     []string{"https://example.com/callback"},
				"smbody":       []string{"Hello, 世界"},
				"smsPointFlag": []string{"1"},
			},
		},
		{
			params: MessageParams{
				Message: Message{
					Dstaddr:  "0987654321",
					Smbody:   "Hello, 世界",
					ClientID: "0aab",
				},
			},
			expected: url.Values{
				"clientid":     []string{"0aab"},
				"dstaddr":      []string{"0987654321"},
				"smbody":       []string{"Hello, 世界"},
				"smsPointFlag": []string{"1"},
			},
		},
		{
			params: MessageParams{
				ObjectID: "batch1",
				Message: Message{
					Dstaddr: "0987654321",
					Smbody:  "Hello, 世界",
				},
			},
			expected: url.Values{
				"dstaddr":      []string{"0987654321"},
				"objectID":     []string{"batch1"},
				"smbody":       []string{"Hello, 世界"},
				"smsPointFlag": []string{"1"},
			},
		},
		{
			params: MessageParams{
				HideDeductedPoints: true,
				Message: Message{
					Dstaddr: "0987654321",
					Smbody:  "Hello, 世界",
				},
			},
			expected: url.Values{
				"dstaddr": []string{"0987654321"},
				"smbody":  []string{"Hello, 世界"},
			},
		},
		{
			params: MessageParams{
				Message: Message{
					Dstaddr: "0987654321",
					Smbody:  "Hello, 世界",
				},
			},
			expected: url.Values{
				"dstaddr":      []string{"0987654321"},
				"smbody":       []string{"Hello, 世界"},
				"smsPointFlag": []string{"1"},
			},
		},
		{
			params: MessageParams{
				ObjectID: "batch1",
				Message: Message{
					Dstaddr:  "0987654321",
					Dlvtime:  "20170101010000",
					Vldtime:  "20170101013000",
					Smbody:   "Hello, 世界",
					Destname: "Bob",
					Response: "https://example.com/callback",
					ClientID: "0aab",
				},
			},
			expected: url.Values{
				"clientid":     []string{"0aab"},
				"destname":     []string{"Bob"},
				"dlvtime":      []string{"20170101010000"},
				"dstaddr":      []string{"0987654321"},
				"objectID":     []string{"batch1"},
				"response":     []string{"https://example.com/callback"},
				"smsPointFlag": []string{"1"},
				"smbody":       []string{"Hello, 世界"},
				"vldtime":      []string{"20170101013000"},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case=%d", i), func(t *testing.T) {
			actual := tc.params.ToData()

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("ToData returned %v, want %v", actual, tc.expected)
			}
		})
	}
}

func TestClient_Send(t *testing.T) {
	testCases := []struct {
		name               string
		params             MessageParams
		response           string
		expectedRequestURI string
		expectedFormData   url.Values
		expectedResponse   *MessageResponse
	}{
		{
			name: "ok",
			params: MessageParams{
				Encoding: defaultEncoding,
				Message: Message{
					Dstaddr: "0987654321",
					Smbody:  "Hello, 世界",
				},
			},
			response: `[1]
msgid=#000000013
statuscode=1
AccountPoint=126
smsPoint=1`,
			expectedRequestURI: "/b2c/mtk/SmSend?CharsetURL=UTF-8",
			expectedFormData: url.Values{
				"password":     []string{"password"},
				"username":     []string{"username"},
				"dstaddr":      []string{"0987654321"},
				"smbody":       []string{"Hello, 世界"},
				"smsPointFlag": []string{"1"},
			},
			expectedResponse: &MessageResponse{
				Results: []*MessageResult{
					{
						Msgid:      "#000000013",
						StatusCode: StatusCode("1"),
						SmsPoint:   Ptr(1),
					},
				},
				AccountPoint: 126,
			},
		},
		{
			name: "default encoding",
			params: MessageParams{
				Message: Message{
					Dstaddr: "0987654321",
					Smbody:  "Hello, 世界",
				},
			},
			response: `[1]
msgid=#000000013
statuscode=1
AccountPoint=126
smsPoint=1`,
			expectedRequestURI: "/b2c/mtk/SmSend?CharsetURL=UTF-8",
			expectedFormData: url.Values{
				"password":     []string{"password"},
				"username":     []string{"username"},
				"dstaddr":      []string{"0987654321"},
				"smbody":       []string{"Hello, 世界"},
				"smsPointFlag": []string{"1"},
			},
			expectedResponse: &MessageResponse{
				Results: []*MessageResult{
					{
						Msgid:      "#000000013",
						StatusCode: StatusCode("1"),
						SmsPoint:   Ptr(1),
					},
				},
				AccountPoint: 126,
			},
		},
		{
			name: "hide deducted points",
			params: MessageParams{
				Message: Message{
					Dstaddr: "0987654321",
					Smbody:  "Hello, 世界",
				},
				HideDeductedPoints: true,
			},
			response: `[1]
msgid=#000000013
statuscode=1
AccountPoint=126`,
			expectedRequestURI: "/b2c/mtk/SmSend?CharsetURL=UTF-8",
			expectedFormData: url.Values{
				"password": []string{"password"},
				"username": []string{"username"},
				"dstaddr":  []string{"0987654321"},
				"smbody":   []string{"Hello, 世界"},
			},
			expectedResponse: &MessageResponse{
				Results: []*MessageResult{
					{
						Msgid:      "#000000013",
						StatusCode: StatusCode("1"),
					},
				},
				AccountPoint: 126,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case=%d %s", i, tc.name), func(t *testing.T) {
			client, mux, teardown := setup()
			defer teardown()

			mux.HandleFunc("/b2c/mtk/SmSend", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				testRequestURI(t, r, tc.expectedRequestURI)
				testFormData(t, r, tc.expectedFormData)
				_, _ = fmt.Fprint(w, tc.response)
			})

			actual, err := client.Send(context.Background(), tc.params)
			if err != nil {
				t.Errorf("Send returned unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, actual) {
				t.Errorf("Send returned %+v, want %+v", actual, tc.expectedResponse)
			}
		})
	}
}

func TestClient_Send_responseStatusCode500(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/b2c/mtk/SmSend", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.Send(context.Background(),
		MessageParams{
			Message: Message{
				Dstaddr: "0987654321",
				Smbody:  "Hello, 世界",
			},
		},
	)

	if err == nil {
		t.Fatal("Send did not return error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Error("error message should contain status code 500")
	}
}

func TestClient_Send_withEmptySmbody(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, err := client.Send(context.Background(),
		MessageParams{
			Message: Message{
				Dstaddr: "0987654321",
				Smbody:  "",
			},
		},
	)

	expectedErr := &ParameterError{Reason: "empty Smbody"}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Send return error %v, want %v", err, expectedErr)
	}
}

func TestBatchMessageOptions_Validate(t *testing.T) {
	testCases := []struct {
		params        BatchMessagesParams
		expectedError error
	}{
		{
			expectedError: &ParameterError{Reason: "empty messages"},
		},
		{
			params: BatchMessagesParams{
				Messages: []Message{
					{
						Dstaddr: "0987654321",
						Smbody:  "Test1",
					},
				},
			},
			expectedError: &ParameterError{Reason: "0: empty ClientID"},
		},
		{
			params: BatchMessagesParams{
				Messages: []Message{
					{
						ClientID: "0aab",
						Smbody:   "Test1",
					},
				},
			},
			expectedError: &ParameterError{Reason: "0: [0aab] empty Dstaddr"},
		},
		{
			params: BatchMessagesParams{
				Messages: []Message{
					{
						ClientID: "0aab",
						Dstaddr:  "0987654321",
						Smbody:   "Test1",
					},
					{
						ClientID: "1aab",
						Dstaddr:  "0987654321",
					},
				},
			},
			expectedError: &ParameterError{Reason: "1: [1aab] empty Smbody"},
		},
		{
			params: BatchMessagesParams{
				Messages: []Message{
					{
						ClientID: "0aab",
						Dstaddr:  "0987654321",
						Smbody:   "Test1",
					},
				},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case=%d", i), func(t *testing.T) {
			err := tc.params.Validate()

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Validate returned error %v, want %v", err, tc.expectedError)
			}
		})
	}
}

func TestBatchMessageOptions_ToData(t *testing.T) {
	testCases := []struct {
		params   BatchMessagesParams
		expected string
	}{
		{
			expected: "",
		},
		{
			params: BatchMessagesParams{
				Messages: []Message{
					{
						ClientID: "0aab",
						Dstaddr:  "0987654321",
						Smbody:   "Test1",
					},
					{
						ClientID: "1aab",
						Dstaddr:  "0987654321",
						Smbody:   "Test2",
					},
				},
			},
			expected: "0aab$$0987654321$$$$$$$$$$Test1\r\n1aab$$0987654321$$$$$$$$$$Test2\r\n",
		},
		{
			params: BatchMessagesParams{
				Messages: []Message{
					{
						ClientID: "0aab",
						Destname: "Bob",
						Dlvtime:  "20170101010000",
						Vldtime:  "20170101012300",
						Dstaddr:  "0987654321",
						Smbody:   "Test1",
						Response: "https://example.com/callback",
					},
				},
			},
			expected: "0aab$$0987654321$$20170101010000$$20170101012300$$Bob$$https://example.com/callback$$Test1\r\n",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case=%d", i), func(t *testing.T) {
			actual := tc.params.ToData()

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("ToData returned %v, want %v", actual, tc.expected)
			}
		})
	}
}

func TestClient_SendBatch(t *testing.T) {
	testCases := []struct {
		name               string
		params             BatchMessagesParams
		response           string
		expectedRequestURI string
		expectedData       string
		expectedResponse   *MessageResponse
	}{
		{
			name: "ok",
			params: BatchMessagesParams{
				Encoding: defaultEncoding,
				Messages: []Message{
					{
						ClientID: "0aab",
						Destname: "Bob",
						Dlvtime:  "20170101010000",
						Vldtime:  "20170101012300",
						Dstaddr:  "0987654321",
						Smbody:   "Test1",
						Response: "https://example.com/callback",
					},
					{
						ClientID: "1aab",
						Destname: "Bob",
						Dstaddr:  "0987654321",
						Smbody:   "Test2",
					},
				},
			},
			response: `[0aab]
msgid=#1010079522
statuscode=1
smsPoint=1
[1aab]
msgid=#1010079523
statuscode=4
smsPoint=1
AccountPoint=98`,
			expectedRequestURI: "/b2c/mtk/SmBulkSend?Encoding_PostIn=UTF-8&password=password&smsPointFlag=1&username=username",
			expectedData:       "0aab$$0987654321$$20170101010000$$20170101012300$$Bob$$https://example.com/callback$$Test1\r\n1aab$$0987654321$$$$$$Bob$$$$Test2\r\n",
			expectedResponse: &MessageResponse{
				Results: []*MessageResult{
					{
						Msgid:      "#1010079522",
						StatusCode: StatusCode("1"),
						SmsPoint:   Ptr(1),
					},
					{
						Msgid:      "#1010079523",
						StatusCode: StatusCode("4"),
						SmsPoint:   Ptr(1),
					},
				},
				AccountPoint: 98,
			},
		},
		{
			name: "default encoding",
			params: BatchMessagesParams{
				Messages: []Message{
					{
						ClientID: "0aab",
						Dstaddr:  "0987654321",
						Smbody:   "Test1",
					},
				},
			},
			response: `[0aab]
msgid=#1010079522
statuscode=1
smsPoint=1
AccountPoint=99`,
			expectedRequestURI: "/b2c/mtk/SmBulkSend?Encoding_PostIn=UTF-8&password=password&smsPointFlag=1&username=username",
			expectedData:       "0aab$$0987654321$$$$$$$$$$Test1\r\n",
			expectedResponse: &MessageResponse{
				Results: []*MessageResult{
					{
						Msgid:      "#1010079522",
						StatusCode: StatusCode("1"),
						SmsPoint:   Ptr(1),
					},
				},
				AccountPoint: 99,
			},
		},
		{
			name: "hide deducted points",
			params: BatchMessagesParams{
				Messages: []Message{
					{
						ClientID: "0aab",
						Dstaddr:  "0987654321",
						Smbody:   "Test1",
					},
				},
				HideDeductedPoints: true,
			},
			response: `[0aab]
msgid=#1010079522
statuscode=1
AccountPoint=99`,
			expectedRequestURI: "/b2c/mtk/SmBulkSend?Encoding_PostIn=UTF-8&password=password&username=username",
			expectedData:       "0aab$$0987654321$$$$$$$$$$Test1\r\n",
			expectedResponse: &MessageResponse{
				Results: []*MessageResult{
					{
						Msgid:      "#1010079522",
						StatusCode: StatusCode("1"),
						SmsPoint:   Ptr(1),
					},
				},
				AccountPoint: 99,
			},
		},
		{
			name: "objectid",
			params: BatchMessagesParams{
				ObjectID: "batch1",
				Messages: []Message{
					{
						ClientID: "0aab",
						Dstaddr:  "0987654321",
						Smbody:   "Test1",
					},
				},
			},
			response: `[0aab]
msgid=#1010079522
statuscode=1
smsPoint=1
AccountPoint=99`,
			expectedRequestURI: "/b2c/mtk/SmBulkSend?Encoding_PostIn=UTF-8&objectID=batch1&password=password&smsPointFlag=1&username=username",
			expectedData:       "0aab$$0987654321$$$$$$$$$$Test1\r\n",
			expectedResponse: &MessageResponse{
				Results: []*MessageResult{
					{
						Msgid:      "#1010079522",
						StatusCode: StatusCode("1"),
						SmsPoint:   Ptr(1),
					},
				},
				AccountPoint: 99,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case=%d %s", i, tc.name), func(t *testing.T) {
			client, mux, teardown := setup()
			defer teardown()

			mux.HandleFunc("/b2c/mtk/SmBulkSend", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				testRequestURI(t, r, tc.expectedRequestURI)
				testData(t, r, tc.expectedData)
				_, _ = fmt.Fprint(w, tc.response)
			})

			actual, err := client.SendBatch(context.Background(), tc.params)
			if err != nil {
				t.Errorf("SendBatch returned unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, actual) {
				t.Errorf("SendBatch returned %+v, want %+v", actual, tc.expectedResponse)
			}
		})
	}
}

func TestClient_SendBatch_responseStatusCode500(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/b2c/mtk/SmBulkSend", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.SendBatch(context.Background(),
		BatchMessagesParams{
			HideDeductedPoints: true,
			Messages: []Message{
				{
					ClientID: "0aab",
					Dstaddr:  "0987654321",
					Smbody:   "Test1",
				},
			},
		},
	)

	if err == nil {
		t.Fatal("Send did not return error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Error("error message should contain status code 500")
	}
}

func TestClient_SendBatch_withEmptySmbody(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/b2c/mtk/SmBulkSend", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURI(t, r, "/b2c/mtk/SmBulkSend?Encoding_PostIn=UTF-8&password=password&username=username")
		testData(t, r, "0aab$$0987654321$$$$$$$$$$Test1\r\n")
		_, _ = fmt.Fprint(w, `[0aab]
msgid=#1010079522
statuscode=1
AccountPoint=99`)
	})

	_, err := client.SendBatch(context.Background(),
		BatchMessagesParams{
			HideDeductedPoints: true,
			Messages: []Message{
				{
					ClientID: "0aab",
					Dstaddr:  "0987654321",
					Smbody:   "",
				},
			},
		},
	)

	expectedErr := &ParameterError{Reason: "0: [0aab] empty Smbody"}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Send return error %v, want %v", err, expectedErr)
	}
}

func Test_parseMessageResponse(t *testing.T) {
	testCases := []struct {
		body             io.Reader
		expectedResponse *MessageResponse
		expectedErr      error
	}{
		{
			body:        strings.NewReader(""),
			expectedErr: &UnexpectedResponseError{Reason: "invalid response"},
		},
		{
			body:        strings.NewReader("foo"),
			expectedErr: &UnexpectedResponseError{Reason: "no clientid"},
		},
		{
			body: strings.NewReader(`[foo]
bar`),
			expectedErr: &UnexpectedResponseError{Reason: "invalid key value pair"},
		},
		{
			body: strings.NewReader("[foo]"),
			expectedResponse: &MessageResponse{
				Results: []*MessageResult{
					{},
				},
			},
		},
		{
			body: strings.NewReader(`[0]
msgid=#000000333
statuscode=0
AccountPoint=92
Duplicate=Y
smsPoint=1
`),
			expectedResponse: &MessageResponse{
				Results: []*MessageResult{
					{
						Msgid:      "#000000333",
						StatusCode: StatusCode("0"),
						SmsPoint:   Ptr(1),
					},
				},
				AccountPoint: 92,
				Duplicate:    Ptr("Y"),
			},
		},
		{
			body: strings.NewReader(`[0]
msgid=#000000333
statuscode=0
[1]
msgid=#000000334
statuscode=1
AccountPoint=92`),
			expectedResponse: &MessageResponse{
				Results: []*MessageResult{
					{
						Msgid:      "#000000333",
						StatusCode: StatusCode("0"),
					},
					{
						Msgid:      "#000000334",
						StatusCode: StatusCode("1"),
					},
				},
				AccountPoint: 92,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case=%d", i), func(t *testing.T) {
			actual, err := parseMessageResponse(tc.body)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("parseMessageResponse returned error %v, want %v", err, tc.expectedErr)
			}
			if !reflect.DeepEqual(actual, tc.expectedResponse) {
				t.Errorf("parseMessageResponse returned %+v, want %+v", actual, tc.expectedResponse)
			}

		})
	}
}

func TestClient_QueryMessageStatus(t *testing.T) {
	testCases := []struct {
		name               string
		params             MessageStatusParams
		response           string
		expectedRequestURI string
		expectedData       string
		expectedResponse   *MessageStatusResponse
	}{
		{
			name:   "ok",
			params: MessageStatusParams{MessageIDs: []string{"1010079522", "1010079523"}},
			response: `1010079522	1	20170101010010	1
1010079523	4	20170101010011	1`,
			expectedRequestURI: "/b2c/mtk/SmQuery?msgid=1010079522%2C1010079523&smsPointFlag=1",
			expectedData:       "1010079522,1010079523",
			expectedResponse: &MessageStatusResponse{
				Statuses: []*MessageStatus{
					{
						MessageResult: MessageResult{
							Msgid:      "1010079522",
							StatusCode: StatusCode("1"),
							SmsPoint:   Ptr(1),
						},
						StatusTime: "20170101010010",
					},
					{
						MessageResult: MessageResult{
							Msgid:      "1010079523",
							StatusCode: StatusCode("4"),
							SmsPoint:   Ptr(1),
						},
						StatusTime: "20170101010011",
					},
				},
			},
		},
		{
			name:               "hide deducted points",
			params:             MessageStatusParams{MessageIDs: []string{"1010079522"}, HideDeductedPoints: true},
			response:           `1010079522	1	20170101010010`,
			expectedRequestURI: "/b2c/mtk/SmQuery?msgid=1010079522",
			expectedData:       "1010079522,1010079523",
			expectedResponse: &MessageStatusResponse{
				Statuses: []*MessageStatus{
					{
						MessageResult: MessageResult{
							Msgid:      "1010079522",
							StatusCode: StatusCode("1"),
						},
						StatusTime: "20170101010010",
					},
				},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case=%d %s", i, tc.name), func(t *testing.T) {
			client, mux, teardown := setup()
			defer teardown()

			mux.HandleFunc("/b2c/mtk/SmQuery", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testRequestURI(t, r, tc.expectedRequestURI)
				_, _ = fmt.Fprint(w, tc.response)
			})

			resp, err := client.QueryMessageStatus(context.Background(), tc.params)
			if err != nil {
				t.Errorf("QueryMessageStatus returned unexpected error: %v", err)
			}

			if !reflect.DeepEqual(resp, tc.expectedResponse) {
				t.Errorf("QueryMessageStatus returned %+v, want %+v", resp, tc.expectedResponse)
			}
		})
	}
}

func TestClient_QueryAccountPoint(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/b2c/mtk/SmQuery", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `AccountPoint=100`)
	})

	ap, err := client.QueryAccountPoint(context.Background())
	if err != nil {
		t.Errorf("QueryAccountPoint returned unexpected error: %v", err)
	}
	if ap != 100 {
		t.Errorf("QueryAccountPoint returned %+v, want %+v", ap, 100)
	}
}

func TestClient_CancelScheduledMessages(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/b2c/mtk/SmCancel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `1010079522=8
1010079523=9`)
	})

	resp, err := client.CancelScheduledMessages(context.Background(), []string{"1010079522", "1010079523"})
	if err != nil {
		t.Errorf("CancelScheduledMessages returned unexpected error: %v", err)
	}

	want := []*CanceledMessage{
		{
			Msgid:      "1010079522",
			StatusCode: StatusCode("8"),
		},
		{
			Msgid:      "1010079523",
			StatusCode: StatusCode("9"),
		},
	}
	if !reflect.DeepEqual(resp, want) {
		t.Errorf("CancelScheduledMessages returned %+v, want %+v", resp, want)
	}
}
