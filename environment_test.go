package scalr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentsList(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	envl, err := client.Environments.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	totalCount := envl.TotalCount
	envTest1, envTest1Cleanup := createEnvironment(t, client)
	defer envTest1Cleanup()

	t.Run("with no list options", func(t *testing.T) {
		envl, err := client.Environments.List(ctx)
		envlIDs := make([]string, len(envl.Items))
		for _, env := range envl.Items {
			envlIDs = append(envlIDs, env.ID)
		}
		require.NoError(t, err)
		assert.Contains(t, envlIDs, envTest1.ID)

		assert.Equal(t, 1, envl.CurrentPage)
		assert.Equal(t, 1+totalCount, envl.TotalCount)
	})

}

func TestEnvironmentsCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	t.Run("when no name is provided", func(t *testing.T) {
		_, err := client.Environments.Create(ctx, EnvironmentCreateOptions{
			Account: &Account{ID: defaultAccountID},
		})
		assert.EqualError(t, err, "name is required")
	})
	t.Run("when no account is provided", func(t *testing.T) {
		_, err := client.Environments.Create(ctx, EnvironmentCreateOptions{
			Name: String(randomString(t)),
		})
		assert.EqualError(t, err, "account is required")
	})
	t.Run("with invalid accountID", func(t *testing.T) {
		env, err := client.Environments.Create(ctx, EnvironmentCreateOptions{
			Account: &Account{ID: badIdentifier},
			Name:    String(randomString(t)),
		})
		assert.Nil(t, env)
		assert.EqualError(t, err, "invalid value for account ID")
	})
	t.Run("with valid options", func(t *testing.T) {
		options := EnvironmentCreateOptions{
			Name:    String("tst-" + randomString(t)),
			Account: &Account{ID: defaultAccountID},
		}

		env, err := client.Environments.Create(ctx, options)
		if err != nil {
			t.Fatal(err)
		}
		// Get a refreshed view of the environment
		_, err = client.Environments.Read(ctx, env.ID)
		require.NoError(t, err)

		defer client.Environments.Delete(ctx, env.ID)

		assert.Equal(t, *options.Name, env.Name)
		assert.Equal(t, *&options.Account.ID, env.Account.ID)
	})

}

func TestEnvironmentsRead(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	envTest, envTestCleanup := createEnvironment(t, client)
	defer envTestCleanup()
	t.Run("when the env exists", func(t *testing.T) {
		_, err := client.Environments.Read(ctx, envTest.ID)
		require.NoError(t, err)
	})

	t.Run("when the env does not exist", func(t *testing.T) {
		_, err := client.Environments.Read(ctx, "notexisting")
		assert.Equal(t, err, ErrResourceNotFound)
	})

	t.Run("with invalid env ID", func(t *testing.T) {
		r, err := client.Environments.Read(ctx, badIdentifier)
		assert.Nil(t, r)
		assert.EqualError(t, err, "invalid value for environment ID")
	})
}

func TestEnvironmentsUpdate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("with valid options", func(t *testing.T) {
		envTest, envTestCleanup := createEnvironment(t, client)

		options := EnvironmentUpdateOptions{
			Name:                  String("tst-" + randomString(t)),
			CostEstimationEnabled: Bool(false),
		}

		env, err := client.Environments.Update(ctx, envTest.ID, options)
		if err != nil {
			envTestCleanup()
		}
		require.NoError(t, err)

		// Make sure we clean up the updated env.
		defer client.Environments.Delete(ctx, env.ID)

		// Also get a fresh result from the API to ensure we get the
		// expected values back.
		refreshed, err := client.Environments.Read(ctx, env.ID)
		require.NoError(t, err)

		for _, item := range []*Environment{
			env,
			refreshed,
		} {
			assert.Equal(t, *options.Name, item.Name)
			assert.Equal(t, *options.CostEstimationEnabled, item.CostEstimationEnabled)
		}
	})

	t.Run("when only updating a subset of fields", func(t *testing.T) {
		envTest, envTestCleanup := createEnvironment(t, client)
		defer envTestCleanup()

		env, err := client.Environments.Update(ctx, envTest.ID, EnvironmentUpdateOptions{})
		require.NoError(t, err)
		assert.Equal(t, envTest.Name, env.Name)
	})
}

func TestEnvironmentsDelete(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("with valid options", func(t *testing.T) {
		envTest, _ := createEnvironment(t, client)

		err := client.Environments.Delete(ctx, envTest.ID)
		require.NoError(t, err)

		// Try fetching the env again - it should error.
		_, err = client.Environments.Read(ctx, envTest.ID)
		assert.Equal(t, err, ErrResourceNotFound)
	})

	t.Run("when the env does not exist", func(t *testing.T) {
		err := client.Environments.Delete(ctx, randomString(t))
		assert.Equal(t, err, ErrResourceNotFound)
	})
}
