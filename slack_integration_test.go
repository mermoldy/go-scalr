package scalr

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSlackIntegrationsCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	env1, deleteEnv1 := createEnvironment(t, client)
	defer deleteEnv1()
	var channelId = os.Getenv("SLACK_CHANNEL_ID")
	if len(channelId) == 0 {
		t.Skip("Set `SLACK_CHANNEL_ID` environment variable to run this test. Container should have connection to slack.")
	}
	slackConnection, err := client.SlackIntegrations.GetConnection(ctx, defaultAccountID)
	if err != nil || slackConnection.ID == "" {
		t.Skip("Scalr instance doesn't have working slack connection.")
	}

	t.Run("with valid options", func(t *testing.T) {

		options := SlackIntegrationCreateOptions{
			Name:         String("test-" + randomString(t)),
			Events:       []string{string(RunApprovalRequiredEvent), string(RunSuccessEvent), string(RunErroredEvent)},
			ChannelId:    &channelId,
			Account:      &Account{ID: defaultAccountID},
			Connection:   slackConnection,
			Environments: []*Environment{env1},
		}

		si, err := client.SlackIntegrations.Create(ctx, options)
		require.NoError(t, err)

		refreshed, err := client.SlackIntegrations.Read(ctx, si.ID)
		require.NoError(t, err)

		for _, item := range []*SlackIntegration{
			si,
			refreshed,
		} {
			assert.NotEmpty(t, item.ID)
			assert.Equal(t, *options.Name, item.Name)
			assert.Equal(t, options.Account, item.Account)
			assert.Equal(t, *options.ChannelId, item.ChannelId)
			assert.Equal(t, options.Events, item.Events)
		}

		err = client.SlackIntegrations.Delete(ctx, si.ID)
		require.NoError(t, err)
	})
}

func TestSlackIntegrationsUpdate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	env1, deleteEnv1 := createEnvironment(t, client)
	defer deleteEnv1()
	env2, deleteEnv2 := createEnvironment(t, client)
	defer deleteEnv2()

	slackConnection, err := client.SlackIntegrations.GetConnection(ctx, defaultAccountID)
	if err != nil || slackConnection.ID == "" {
		t.Skip("Scalr instance doesn't have working slack connection.")
	}
	var channelId = os.Getenv("SLACK_CHANNEL_ID")
	if len(channelId) == 0 {
		t.Skip("Set `SLACK_CHANNEL_ID` environment variable to run this test. Container should have connection to slack.")
	}

	si, deleteSlack := createSlackIntegration(t, client, slackConnection, &channelId, env1)
	defer deleteSlack()
	t.Run("with valid options", func(t *testing.T) {

		options := SlackIntegrationUpdateOptions{
			Name:         String("test-" + randomString(t)),
			Events:       []string{RunApprovalRequiredEvent, RunErroredEvent},
			Environments: []*Environment{env2},
		}

		si, err := client.SlackIntegrations.Update(ctx, si.ID, options)
		require.NoError(t, err)

		refreshed, err := client.SlackIntegrations.Read(ctx, si.ID)
		require.NoError(t, err)

		for _, item := range []*SlackIntegration{
			si,
			refreshed,
		} {
			assert.NotEmpty(t, item.ID)
			assert.Equal(t, *options.Name, item.Name)
			assert.Equal(t, options.Events, item.Events)
		}
	})
}

func TestSlackIntegrationsList(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	env1, deleteEnv1 := createEnvironment(t, client)
	defer deleteEnv1()
	env2, deleteEnv2 := createEnvironment(t, client)
	defer deleteEnv2()

	slackConnection, err := client.SlackIntegrations.GetConnection(ctx, defaultAccountID)
	if err != nil || slackConnection.ID == "" {
		t.Skip("Scalr instance doesn't have working slack connection.")
	}
	var channelId = os.Getenv("SLACK_CHANNEL_ID")
	if len(channelId) == 0 {
		t.Skip("Set `SLACK_CHANNEL_ID` environment variable to run this test. Container should have connection to slack.")
	}

	si, deleteSlack := createSlackIntegration(t, client, slackConnection, &channelId, env1)
	defer deleteSlack()
	si2, deleteSlack2 := createSlackIntegration(t, client, slackConnection, &channelId, env2)
	defer deleteSlack2()
	t.Run("with valid options", func(t *testing.T) {

		options := SlackIntegrationListOptions{
			Account: String(defaultAccountID),
		}

		sil, err := client.SlackIntegrations.List(ctx, options)
		require.NoError(t, err)

		assert.Equal(t, 2, sil.TotalCount)
		expectedIDs := []string{si.ID, si2.ID}
		actualIDs := make([]string, len(sil.Items))
		for i, s := range sil.Items {
			actualIDs[i] = s.ID
		}
		assert.ElementsMatch(t, expectedIDs, actualIDs)
	})
}
