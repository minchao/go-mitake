package mitake

import "testing"

func TestStatusCode_String(t *testing.T) {
	actual := StatusServiceError.String()

	if actual != "系統發生錯誤，請聯絡三竹資訊窗口人員" {
		t.Error("StatusServiceError.String() returned unexpected value")
	}
}
