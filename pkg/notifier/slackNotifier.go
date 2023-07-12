package notifier

import (
	"fmt"
	"probe/api/v1alpha1"

	"github.com/slack-go/slack"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Function for sending the Slack notification
func NotifySlack(client client.Client, namespace string, slackDetails v1alpha1.SlackDetails, message string) error {
	token, err := secretResolver(client, namespace, slackDetails.SlackToken)
	if err != nil {
		return err
	}

	channel := slackDetails.SlackChannel
	api := slack.New(token)
	channelID, timestamp, err := api.PostMessage(
		channel,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		return fmt.Errorf("error sending slack message: %s", err)
	}

	_ = channelID // not used, ignore
	_ = timestamp // not used, ignore
	return nil
}
