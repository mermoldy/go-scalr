package scalr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicyGroupEnvironmentsCreate(t *testing.T) {

	client := testClient(t)
	ctx := context.Background()

	envTest, envTestCleanup := createEnvironment(t, client)
	defer envTestCleanup()

	policyGroup, policyGroupCleanup := createPolicyGroup(t, client, nil)
	defer policyGroupCleanup()

	t.Run("with valid options", func(t *testing.T) {
		options := PolicyGroupEnvironmentsCreateOptions{
			PolicyGroupID:           policyGroup.ID,
			PolicyGroupEnvironments: []*PolicyGroupEnvironment{{ID: envTest.ID}},
		}

		err := client.PolicyGroupEnvironments.Create(ctx, options)

		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.PolicyGroups.Read(ctx, policyGroup.ID)
		require.NoError(t, err)

		for _, item := range refreshed.Environments {
			assert.Equal(t, envTest.ID, item.ID)
		}

		func() {
			client.PolicyGroupEnvironments.Delete(
				ctx,
				PolicyGroupEnvironmentDeleteOptions{
					PolicyGroupID: policyGroup.ID,
					EnvironmentID: envTest.ID,
				},
			)
		}()
	})

	t.Run("with empty options", func(t *testing.T) {
		err := client.PolicyGroupEnvironments.Create(ctx, PolicyGroupEnvironmentsCreateOptions{})
		assert.EqualError(t, err, "invalid value for policy group ID")
	})

	t.Run("when options has an invalid environment", func(t *testing.T) {
		var envID = "env-123"
		options := PolicyGroupEnvironmentsCreateOptions{
			PolicyGroupID:           policyGroup.ID,
			PolicyGroupEnvironments: []*PolicyGroupEnvironment{{ID: envID}},
		}

		err := client.PolicyGroupEnvironments.Create(ctx, options)
		assert.NotEmpty(t, err)
	})

}

func TestPolicyGroupEnvironmentDelete(t *testing.T) {

	client := testClient(t)
	ctx := context.Background()

	envTest, envTestCleanup := createEnvironment(t, client)
	policyGroup, policyGroupCleanup := createPolicyGroup(t, client, nil)
	policyGroupEnvironmentLinkCleanup := linkPolicyGroupToEnvironment(t, client, policyGroup, envTest)
	defer policyGroupEnvironmentLinkCleanup()
	defer policyGroupCleanup()
	defer envTestCleanup()

	t.Run("with valid options", func(t *testing.T) {
		err := client.PolicyGroupEnvironments.Delete(ctx, PolicyGroupEnvironmentDeleteOptions{
			PolicyGroupID: policyGroup.ID,
			EnvironmentID: envTest.ID,
		})
		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.PolicyGroups.Read(ctx, policyGroup.ID)
		require.NoError(t, err)
		assert.Empty(t, refreshed.Environments)
	})

	t.Run("without a valid policy group ID", func(t *testing.T) {
		err := client.PolicyGroupEnvironments.Delete(ctx, PolicyGroupEnvironmentDeleteOptions{
			PolicyGroupID: badIdentifier,
			EnvironmentID: envTest.ID,
		})
		assert.EqualError(t, err, "invalid value for policy group ID")
	})
}
