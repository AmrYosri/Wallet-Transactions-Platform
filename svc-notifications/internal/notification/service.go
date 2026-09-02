package notification

import (
	"context"
	"log"
	"time"

)


type Service struct{
	repo *Repository
}

func NewService(repo *Repository) *Service{
	return &Service{repo: repo}
}

func (s *Service) Notify(ctx context.Context,req *NotificationRequest)(*NotificationResponse,error){
	s.processSMS(ctx,req)
	s.processPush(ctx,req)
	return &NotificationResponse{Status: "accepted"},nil
}

func (s *Service) processSMS(ctx context.Context, req *NotificationRequest){
	smsNotification:=&SMSNotification{
		TransactionID: req.TransactionID,
		Phone: req.Phone,
		Message: req.SMSMessage,
		Status: SMSStatusPending,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateSMS(ctx,smsNotification); err != nil{
		log.Printf("failed to create sms notification record: %v",err)
		return
	}
	ok,reason := stubSendSMS(req.Phone,req.SMSMessage)
	if ok{
		if err := s.repo.MarkSMSSent(ctx,smsNotification.ID); err !=nil{
			log.Printf("failed to mark sms sent: %v",err)

		}
		return
	}
	if err := s.repo.MarkSMSFailed(ctx ,smsNotification.ID,reason); err !=nil{
		log.Printf("failed to mark sms failed: %v",err)
	}
}

func (s *Service) processPush(ctx context.Context,req *NotificationRequest){
	pushNotification := &PushNotification{
		TransactionID: req.TransactionID,
		DeviceToken: req.DeviceToken,
		Title: req.PushTitle,
		Body: req.PushBody,
		Status: PushStatusPending,
		CreatedAt: time.Now(),

	}

	if err := s.repo.CreatePush(ctx,pushNotification); err !=nil{
		log.Printf("failed to create push notification record: %v", err)
		return
	}
	ok ,reason := stubSendPush(req.DeviceToken,req.PushTitle,req.PushBody)
	if ok {
		if err := s.repo.MarkPushSent(ctx,pushNotification.ID); err != nil{
			log.Printf("failed to mark push sent: %v", err)
		}
		return
	}

	if err := s.repo.MarkPushFailed(ctx,pushNotification.ID,reason); err != nil{
		log.Printf("failed to mark push failed: %v", err)
	}

}

func stubSendSMS(phone, message string) (success bool, failureReason string) {
	if phone == "" {
		log.Printf("[STUB SMS] send FAILED — missing phone number")
		return false, "stub: missing phone number"
	}
	log.Printf("[STUB SMS] send OK to %s: %q", phone, message)
	return true, ""
}

func stubSendPush(deviceToken,title,body string)(success bool , failureReason string){
	if deviceToken =="" {
		log.Printf("[STUB PUSH] send FAILED - missing device token")
		return false , "stub : missing device token"
	}
	log.Printf("[STUB PUSH] send OK to %s : %q / %q ", deviceToken,title,body)
	return true ,""
}




