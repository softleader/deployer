package slack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v28/github"
	"github.com/nlopes/slack"
	"github.com/softleader/deployer/cmd/docker"
	"github.com/softleader/deployer/models"
	"golang.org/x/oauth2"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrMissingSlackWebhookURL 代表沒有 Webhook URL
	ErrMissingSlackWebhookURL = errors.New(`missing slack webhook URL`)
)

func Post(config models.Config, serviceId, image string) error {
	if config.SlackAPI.WebHookURL == "" {
		return ErrMissingSlackWebhookURL
	}
	tag := after(image, ":")
	payload := &slack.WebhookMessage{
		Text: fmt.Sprintf("SIT %s 過版", between(image, "/", ":")),
	}
	for _, attachment := range newAttachments(config, serviceId, tag) {
		payload.Attachments = append(payload.Attachments, attachment)
	}
	return slack.PostWebhook(config.SlackAPI.WebHookURL, payload)
}

func newAttachments(config models.Config, serviceId, tag string) (attachments []slack.Attachment) {
	release := slack.Attachment{
		Title:  tag,
		Footer: "github.com",
		Ts:     json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	arg, spec, err := docker.ServiceSpec(serviceId)
	if err != nil {
		fmt.Printf("failed inspecting docker service by command [%v]: %v\n", arg, err)
	} else {
		if val, found := spec.Labels["com.docker.stack.namespace"]; found {
			attachments = append(attachments, slack.Attachment{
				Title:     replaceLast(val, "-", "/"),
				TitleLink: fmt.Sprintf("http://softleader.com.tw:5678/services/%v?q=%v", val, serviceId),
				Footer:    "http://softleader.com.tw:5678",
				Ts:        json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
			})
		}
		if val, found := spec.Labels["github"]; found {
			github := strings.Split(val, "/")
			owner := github[0]
			repo := github[1]
			if token, ok := config.GitHubToken[owner]; !ok {
				fmt.Printf("not found any github token for '%v' in config file\n", owner)
			} else {
				if r, _, err := getReleaseByTag(token, owner, repo, tag); err != nil {
					fmt.Printf("failed fetching github [%v/%v] release by tag [%v]: %v\n", owner, repo, tag, err)
				} else {
					release.Title = r.GetTagName()
					release.TitleLink = r.GetHTMLURL()
					release.AuthorName = r.GetAuthor().GetLogin()
					release.AuthorLink = r.GetAuthor().GetHTMLURL()
					release.AuthorIcon = r.GetAuthor().GetAvatarURL()
					release.Footer = fmt.Sprintf("https://github.com/%v/%v", owner, repo)
					release.Ts = json.Number(strconv.FormatInt(r.GetPublishedAt().Unix(), 10))
				}
			}
		}
	}
	attachments = append(attachments, release)
	return
}

func getReleaseByTag(token, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
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

// Get substring after a string.
func after(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

// Get substring after a string.
func replaceLast(value, target, to string) string {
	pos := strings.LastIndex(value, target)
	if pos == -1 {
		return value
	}
	adjustedPos := pos + len(target)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[0:adjustedPos-len(target)] + to + value[adjustedPos:len(value)]
}
