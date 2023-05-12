package scalr

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSlackIntegrationsCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	env1, deleteEnv1 := createEnvironment(t, client)
	defer deleteEnv1()

	slackConnection, err := client.SlackIntegrations.GetConnection(ctx, defaultAccountID)
	if err != nil {
		println(err.Error())
		return
	}
	if slackConnection.ID == "" {
		t.Skip("Scalr instance doesn't have working slack connection.")
	}
	slackChannels, _ := client.SlackIntegrations.GetChannels(ctx, defaultAccountID, SlackChannelListOptions{})
	var channelId string
	for _, channel := range slackChannels.Items {
		channelId = channel.ID
		break
	}
	t.Run("with valid options", func(t *testing.T) {

		options := SlackIntegrationCreateOptions{
			Name:        String("test-" + randomString(t)),
			Events:      []SlackEvent{RunApprovalRequiredEvent, RunSuccessEvent, RunErroredEvent},
			ChannelId:   &channelId,
			Account:     &Account{ID: defaultAccountID},
			Connection:  slackConnection,
			Environment: env1,
		}

		_, err := client.SlackIntegrations.Create(ctx, options)
		require.NoError(t, err)
	})
}
