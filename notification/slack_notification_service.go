package notification

import (
	"os"
	"time"
	"errors"
	"strconv"
	"encoding/json"
	"github.com/slack-go/slack"
	"database-backup-manager/constants"
)

type slackNotificationService struct {
	webhookUrl string
}

func NewSlackNotificationService() (*slackNotificationService, error) {
	sns := new(slackNotificationService)
	sns.webhookUrl = os.Getenv(constants.SLACK_WEBHOOK_URL)
	if (len(sns.webhookUrl) < 16) {
		return nil, errors.New("Webhook Url not set!")
	}
	return sns, nil
}

func (ns *slackNotificationService) Send(message string) error {
	attachment := slack.Attachment{
		Text:          message,
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	return slack.PostWebhook(ns.webhookUrl, &msg)
}