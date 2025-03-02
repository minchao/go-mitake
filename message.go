package mitake

// StatusCode of Mitake API.
type StatusCode string

func (c StatusCode) String() string {
	return statusCodeMap[c]
}

// List of Mitake API status codes.
const (
	StatusServiceError                    = StatusCode("*")
	StatusSMSTemporarilyUnavailable       = StatusCode("a")
	StatusSMSTemporarilyUnavailableB      = StatusCode("b")
	StatusUsernameRequired                = StatusCode("c")
	StatusPasswordRequired                = StatusCode("d")
	StatusUsernameOrPasswordError         = StatusCode("e")
	StatusAccountExpired                  = StatusCode("f")
	StatusAccountDisabled                 = StatusCode("h")
	StatusInvalidConnectionAddress        = StatusCode("k")
	StatusReachedMaxConcurrentConnections = StatusCode("l")
	StatusChangePasswordRequired          = StatusCode("m")
	StatusPasswordExpired                 = StatusCode("n")
	StatusPermissionDenied                = StatusCode("p")
	StatusServiceTemporarilyUnavailable   = StatusCode("r")
	StatusAccountingFailure               = StatusCode("s")
	StatusSMSExpired                      = StatusCode("t")
	StatusSMSBodyEmpty                    = StatusCode("u")
	StatusInvalidPhoneNumber              = StatusCode("v")
	StatusQueryCountExceedsLimit          = StatusCode("w")
	StatusFileToLargeToSend               = StatusCode("x")
	StatusInvalidParameter                = StatusCode("y")
	StatusNoDataFound                     = StatusCode("z")

	StatusReservationForDelivery = StatusCode("0")
	StatusCarrierAccepted        = StatusCode("1")
	StatusCarrierAccepted2       = StatusCode("2")
	StatusCarrierAccepted3       = StatusCode("3")
	StatusDelivered              = StatusCode("4")
	StatusContentError           = StatusCode("5")
	StatusPhoneNumberError       = StatusCode("6")
	StatusSMSDisable             = StatusCode("7")
	StatusDeliveryTimeout        = StatusCode("8")
	StatusReservationCanceled    = StatusCode("9")
)

var statusCodeMap = map[StatusCode]string{
	StatusServiceError:                    "系統發生錯誤，請聯絡三竹資訊窗口人員",
	StatusSMSTemporarilyUnavailable:       "簡訊發送功能暫時停止服務，請稍候再試",
	StatusSMSTemporarilyUnavailableB:      "簡訊發送功能暫時停止服務，請稍候再試",
	StatusUsernameRequired:                "請輸入帳號",
	StatusPasswordRequired:                "請輸入密碼",
	StatusUsernameOrPasswordError:         "帳號、密碼錯誤",
	StatusAccountExpired:                  "帳號已過期",
	StatusAccountDisabled:                 "帳號已被停用",
	StatusInvalidConnectionAddress:        "無效的連線位址",
	StatusReachedMaxConcurrentConnections: "帳號已達到同時連線數上限",
	StatusChangePasswordRequired:          "必須變更密碼，在變更密碼前，無法使用簡訊發送服務",
	StatusPasswordExpired:                 "密碼已逾期，在變更密碼前，將無法使用簡訊發送服務",
	StatusPermissionDenied:                "沒有權限使用外部Http程式",
	StatusServiceTemporarilyUnavailable:   "系統暫停服務，請稍後再試",
	StatusAccountingFailure:               "帳務處理失敗，無法發送簡訊",
	StatusSMSExpired:                      "簡訊已過期",
	StatusSMSBodyEmpty:                    "簡訊內容不得為空白",
	StatusInvalidPhoneNumber:              "無效的手機號碼",
	StatusQueryCountExceedsLimit:          "查詢筆數超過上限",
	StatusFileToLargeToSend:               "發送檔案過大，無法發送簡訊",
	StatusInvalidParameter:                "參數錯誤",
	StatusNoDataFound:                     "查無資料",

	StatusReservationForDelivery: "預約傳送中",
	StatusCarrierAccepted:        "已送達業者",
	StatusCarrierAccepted2:       "已送達業者",
	StatusCarrierAccepted3:       "已送達業者",
	StatusDelivered:              "已送達手機",
	StatusContentError:           "內容有錯誤",
	StatusPhoneNumberError:       "門號有錯誤",
	StatusSMSDisable:             "簡訊已停用",
	StatusDeliveryTimeout:        "逾時無送達",
	StatusReservationCanceled:    "預約已取消",
}

type Message struct {
	ClientID string // A unique identifier from client to identify SMS message
	Dstaddr  string // Required, Destination phone number
	Smbody   string // Required, The text of the message you want to send, use ASCII code 6 to represent a new line
	Dlvtime  string // Scheduled delivery time, format: YYYYMMDDHHMMSS
	Vldtime  string // Validity period, format: YYYYMMDDHHMMSS
	Destname string // Destination receiver name
	Response string // Callback URL to receive the delivery receipt of the message
}
