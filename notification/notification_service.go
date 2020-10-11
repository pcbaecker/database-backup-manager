package notification

import (
	"os"
	"fmt"
	"errors"
	"database-backup-manager/constants"
)

type NotificationService interface {
	Send(message string) error
}

func NewNotificationService() (NotificationService, error) {
	notification_type := os.Getenv(constants.NOTIFICATION_TYPE)
	if (notification_type == "slack") {
		return NewSlackNotificationService()
	} else if (len(notification_type) > 0) {
		return nil, errors.New("Unspecified notification type = " + notification_type)
	}

	return newDummyNotificationService()
}

func newDummyNotificationService() (NotificationService, error) {
	dns := new(dummyNotificationService)
	return dns, nil
}

type dummyNotificationService struct {

}

func (dns *dummyNotificationService) Send(message string) error {
	fmt.Println("No notification service specified!")
	return nil
}