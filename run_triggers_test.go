package scalr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunTriggersCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	env1Test, env1TestCleanup := createEnvironment(t, client)
	defer env1TestCleanup()

	env2Test, envTest2Cleanup := createEnvironment(t, client)
	defer envTest2Cleanup()

	wsEnv1Test1, wsEnv1Test1Cleanup := createWorkspace(t, client, env1Test)
	defer wsEnv1Test1Cleanup()
	wsEnv1Test2, wsEnv1Test2Cleanup := createWorkspace(t, client, env1Test)
	defer wsEnv1Test2Cleanup()
	wsEnv2Test1, wsEnv2Test1Cleanup := createWorkspace(t, client, env2Test)
	defer wsEnv2Test1Cleanup()

	t.Run("missing downstream workspace", func(t *testing.T) {
		options := RunTriggerCreateOptions{
			Upstream: &Upstream{ID: wsEnv1Test1.ID},
		}
		trigger, err := client.RunTriggers.Create(ctx, options)
		require.Error(t, err)
		assert.EqualError(t, err, "downstream ID is required")
		assert.Nil(t, trigger)
	})

	t.Run("missing upstream workspace", func(t *testing.T) {
		options := RunTriggerCreateOptions{
			Downstream: &Downstream{ID: wsEnv1Test1.ID},
		}
		trigger, err := client.RunTriggers.Create(ctx, options)
		require.Error(t, err)
		assert.EqualError(t, err, "upstream ID is required")
		assert.Nil(t, trigger)
	})

	t.Run("from different environments", func(t *testing.T) {
		options := RunTriggerCreateOptions{
			Downstream: &Downstream{ID: wsEnv1Test1.ID},
			Upstream:   &Upstream{ID: wsEnv2Test1.ID},
		}
		trigger, err := client.RunTriggers.Create(ctx, options)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "The downstream and upstream workspaces must be within the same Scalr environment")
		assert.Nil(t, trigger)
	})

	t.Run("create trigger", func(t *testing.T) {
		options := RunTriggerCreateOptions{
			Downstream: &Downstream{ID: wsEnv1Test1.ID},
			Upstream:   &Upstream{ID: wsEnv1Test2.ID},
		}
		trigger, err := client.RunTriggers.Create(ctx, options)
		require.NoError(t, err)
		assert.NotEmpty(t, trigger.ID)
		assert.NotEmpty(t, trigger.CreatedAt)
		assert.Equal(t, wsEnv1Test1.ID, trigger.Downstream.ID)
		assert.Equal(t, wsEnv1Test2.ID, trigger.Upstream.ID)
	})

}
