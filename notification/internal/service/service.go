package service

import "fmt"

type NotifService struct {
}

func NewNotifService() *NotifService {
	return &NotifService{}
}

func (n *NotifService) SendItToTheEmail(email string) (string, error) {
	//
	resp := fmt.Sprintf("msg has been sent to email: %s", email)
	return resp, nil
}
func (n *NotifService) SentItToThePhone(phone string) (string, error) {
	resp := fmt.Sprintf("msg has been sent to phone: %s", phone)
	return resp, nil
}
