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

	channelID := slackDetails.SlackChannelID
	api := slack.New(token)
	resChannel, respTimestamp, err := api.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		return fmt.Errorf("error sending slack message: %s", err)
	}

	_ = resChannel    // not used, ignore
	_ = respTimestamp // not used, ignore
	return nil
}
