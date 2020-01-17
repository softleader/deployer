package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/softleader/deployer/models"
	"strings"
)

var (
	// ErrMissingSlackWebhookURL 代表沒有 Webhook URL
	ErrMissingSlackWebhookURL = errors.New(`missing slack webhook URL`)
)

func Post(api models.SlackAPI, image, title, titleLink, authorLink, authorName, authorIcon, ts string) error {
	if api.WebHookURL == "" {
		return ErrMissingSlackWebhookURL
	}
	payload := &slack.WebhookMessage{
		Text: fmt.Sprintf("SIT %s 過版", between(image, "/", ":")),
	}
	attachment := slack.Attachment{
		Title:      title,
		TitleLink:  titleLink,
		AuthorName: authorName,
		AuthorLink: authorLink,
		AuthorIcon: authorIcon,
		Footer:     api.Footer,
		Ts:         json.Number(ts),
	}
	payload.Attachments = append(payload.Attachments, attachment)
	return slack.PostWebhook(api.WebHookURL, payload)
}

// Get substring between two strings
func between(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}
