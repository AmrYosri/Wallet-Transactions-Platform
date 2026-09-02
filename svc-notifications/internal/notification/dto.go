package notification


type NotificationRequest struct {
	TransactionID string `json:"transaction_id"`
	Phone string `json:"phone"`
	DeviceToken string `json:"device_token"`

	SMSMessage string `json:"sms_message"`
	PushTitle string `json:"push_title"`
	PushBody string `json:"push_body"`


}

type NotificationResponse struct {
	Status string `json:"status"`
}