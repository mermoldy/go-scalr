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

	orgTest1, orgTest1Cleanup := createEnvironment(t, client)
	defer orgTest1Cleanup()
	orgTest2, orgTest2Cleanup := createEnvironment(t, client)
	defer orgTest2Cleanup()

	t.Run("with no list options", func(t *testing.T) {
		orgl, err := client.Environments.List(ctx, EnvironmentListOptions{})
		require.NoError(t, err)
		assert.Contains(t, orgl.Items, orgTest1)
		assert.Contains(t, orgl.Items, orgTest2)

		t.Skip("paging not supported yet in API")
		assert.Equal(t, 1, orgl.CurrentPage)
		assert.Equal(t, 2, orgl.TotalCount)
	})

	t.Run("with list options", func(t *testing.T) {
		t.Skip("paging not supported yet in API")
		// Request a page number which is out of range. The result should
		// be successful, but return no results if the paging options are
		// properly passed along.
		orgl, err := client.Environments.List(ctx, EnvironmentListOptions{
			ListOptions: ListOptions{
				PageNumber: 999,
				PageSize:   100,
			},
		})
		require.NoError(t, err)
		assert.Empty(t, orgl)
		assert.Equal(t, 999, orgl.CurrentPage)
		assert.Equal(t, 2, orgl.TotalCount)
	})
}

func TestEnvironmentsCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("with valid options", func(t *testing.T) {
		options := EnvironmentCreateOptions{
			Name:  String(randomString(t)),
			Email: String(randomString(t) + "@scalr.local"),
		}

		org, err := client.Environments.Create(ctx, options)
		require.NoError(t, err)

		// Make sure we clean up the created org.
		defer client.Environments.Delete(ctx, org.Name)

		assert.Equal(t, *options.Name, org.Name)
		assert.Equal(t, *options.Email, org.Email)
	})

	t.Run("when no email is provided", func(t *testing.T) {
		org, err := client.Environments.Create(ctx, EnvironmentCreateOptions{
			Name: String("foo"),
		})
		assert.Nil(t, org)
		assert.EqualError(t, err, "email is required")
	})

	t.Run("when no name is provided", func(t *testing.T) {
		_, err := client.Environments.Create(ctx, EnvironmentCreateOptions{
			Email: String("foo@bar.com"),
		})
		assert.EqualError(t, err, "name is required")
	})

	t.Run("with invalid name", func(t *testing.T) {
		org, err := client.Environments.Create(ctx, EnvironmentCreateOptions{
			Name:  String(badIdentifier),
			Email: String("foo@bar.com"),
		})
		assert.Nil(t, org)
		assert.EqualError(t, err, "invalid value for name")
	})
}

func TestEnvironmentsRead(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	orgTest, orgTestCleanup := createEnvironment(t, client)
	defer orgTestCleanup()

	t.Run("when the org exists", func(t *testing.T) {
		org, err := client.Environments.Read(ctx, orgTest.Name)
		require.NoError(t, err)
		assert.Equal(t, orgTest, org)

		t.Run("permissions are properly decoded", func(t *testing.T) {
			assert.True(t, org.Permissions.CanDestroy)
		})

		t.Run("timestamps are populated", func(t *testing.T) {
			assert.NotEmpty(t, org.CreatedAt)
			assert.NotEmpty(t, org.TrialExpiresAt)
		})
	})

	t.Run("with invalid name", func(t *testing.T) {
		org, err := client.Environments.Read(ctx, badIdentifier)
		assert.Nil(t, org)
		assert.EqualError(t, err, "invalid value for environment")
	})

	t.Run("when the org does not exist", func(t *testing.T) {
		_, err := client.Environments.Read(ctx, randomString(t))
		assert.Error(t, err)
	})
}

func TestEnvironmentsUpdate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("with valid options", func(t *testing.T) {
		orgTest, orgTestCleanup := createEnvironment(t, client)

		options := EnvironmentUpdateOptions{
			Name:            String(randomString(t)),
			Email:           String(randomString(t) + "@scalr.local"),
			SessionTimeout:  Int(3600),
			SessionRemember: Int(3600),
		}

		org, err := client.Environments.Update(ctx, orgTest.Name, options)
		if err != nil {
			orgTestCleanup()
		}
		require.NoError(t, err)

		// Make sure we clean up the renamed org.
		defer client.Environments.Delete(ctx, org.Name)

		// Also get a fresh result from the API to ensure we get the
		// expected values back.
		refreshed, err := client.Environments.Read(ctx, *options.Name)
		require.NoError(t, err)

		for _, item := range []*Environment{
			org,
			refreshed,
		} {
			assert.Equal(t, *options.Name, item.Name)
			assert.Equal(t, *options.Email, item.Email)
			assert.Equal(t, *options.SessionTimeout, item.SessionTimeout)
			assert.Equal(t, *options.SessionRemember, item.SessionRemember)
		}
	})

	t.Run("with invalid name", func(t *testing.T) {
		org, err := client.Environments.Update(ctx, badIdentifier, EnvironmentUpdateOptions{})
		assert.Nil(t, org)
		assert.EqualError(t, err, "invalid value for environment")
	})

	t.Run("when only updating a subset of fields", func(t *testing.T) {
		orgTest, orgTestCleanup := createEnvironment(t, client)
		defer orgTestCleanup()

		org, err := client.Environments.Update(ctx, orgTest.Name, EnvironmentUpdateOptions{})
		require.NoError(t, err)
		assert.Equal(t, orgTest.Name, org.Name)
		assert.Equal(t, orgTest.Email, org.Email)
	})
}

func TestEnvironmentsDelete(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("with valid options", func(t *testing.T) {
		orgTest, _ := createEnvironment(t, client)

		err := client.Environments.Delete(ctx, orgTest.Name)
		require.NoError(t, err)

		// Try fetching the org again - it should error.
		_, err = client.Environments.Read(ctx, orgTest.Name)
		assert.Equal(t, err, ErrResourceNotFound)
	})

	t.Run("with invalid name", func(t *testing.T) {
		err := client.Environments.Delete(ctx, badIdentifier)
		assert.EqualError(t, err, "invalid value for environment")
	})
}

func TestEnvironmentsCapacity(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	orgTest, orgTestCleanup := createEnvironment(t, client)
	defer orgTestCleanup()

	wTest1, _ := createWorkspace(t, client, orgTest)
	wTest2, _ := createWorkspace(t, client, orgTest)
	wTest3, _ := createWorkspace(t, client, orgTest)
	wTest4, _ := createWorkspace(t, client, orgTest)

	t.Run("without queued runs", func(t *testing.T) {
		c, err := client.Environments.Capacity(ctx, orgTest.Name)
		require.NoError(t, err)
		assert.Equal(t, 0, c.Pending)
		assert.Equal(t, 0, c.Running)
	})

	// For this test FRQ should be enabled and have a
	// limit of 2 concurrent runs per environment.
	t.Run("with queued runs", func(t *testing.T) {
		_, _ = createRun(t, client, wTest1)
		_, _ = createRun(t, client, wTest2)
		_, _ = createRun(t, client, wTest3)
		_, _ = createRun(t, client, wTest4)

		c, err := client.Environments.Capacity(ctx, orgTest.Name)
		require.NoError(t, err)
		assert.Equal(t, 2, c.Pending)
		assert.Equal(t, 2, c.Running)
	})

	t.Run("with invalid name", func(t *testing.T) {
		org, err := client.Environments.Read(ctx, badIdentifier)
		assert.Nil(t, org)
		assert.EqualError(t, err, "invalid value for environment")
	})

	t.Run("when the org does not exist", func(t *testing.T) {
		_, err := client.Environments.Read(ctx, randomString(t))
		assert.Error(t, err)
	})
}

func TestEnvironmentsEntitlements(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	orgTest, orgTestCleanup := createEnvironment(t, client)
	defer orgTestCleanup()

	t.Run("when the org exists", func(t *testing.T) {
		entitlements, err := client.Environments.Entitlements(ctx, orgTest.Name)
		require.NoError(t, err)

		assert.NotEmpty(t, entitlements.ID)
		assert.True(t, entitlements.Operations)
		assert.True(t, entitlements.PrivateModuleRegistry)
		assert.True(t, entitlements.Sentinel)
		assert.True(t, entitlements.StateStorage)
		assert.True(t, entitlements.Teams)
		assert.True(t, entitlements.VCSIntegrations)
	})

	t.Run("with invalid name", func(t *testing.T) {
		entitlements, err := client.Environments.Entitlements(ctx, badIdentifier)
		assert.Nil(t, entitlements)
		assert.EqualError(t, err, "invalid value for environment")
	})

	t.Run("when the org does not exist", func(t *testing.T) {
		_, err := client.Environments.Entitlements(ctx, randomString(t))
		assert.Equal(t, ErrResourceNotFound, err)
	})
}

func TestEnvironmentsRunQueue(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	orgTest, orgTestCleanup := createEnvironment(t, client)
	defer orgTestCleanup()

	wTest1, _ := createWorkspace(t, client, orgTest)
	wTest2, _ := createWorkspace(t, client, orgTest)
	wTest3, _ := createWorkspace(t, client, orgTest)
	wTest4, _ := createWorkspace(t, client, orgTest)

	t.Run("without queued runs", func(t *testing.T) {
		rq, err := client.Environments.RunQueue(ctx, orgTest.Name, RunQueueOptions{})
		require.NoError(t, err)
		assert.Equal(t, 0, len(rq.Items))
	})

	// Create a couple or runs to fill the queue.
	rTest1, _ := createRun(t, client, wTest1)
	rTest2, _ := createRun(t, client, wTest2)
	rTest3, _ := createRun(t, client, wTest3)
	rTest4, _ := createRun(t, client, wTest4)

	// For this test FRQ should be enabled and have a
	// limit of 2 concurrent runs per environment.
	t.Run("with queued runs", func(t *testing.T) {
		rq, err := client.Environments.RunQueue(ctx, orgTest.Name, RunQueueOptions{})
		require.NoError(t, err)

		found := []string{}
		for _, r := range rq.Items {
			found = append(found, r.ID)
		}

		assert.Contains(t, found, rTest1.ID)
		assert.Contains(t, found, rTest2.ID)
		assert.Contains(t, found, rTest3.ID)
		assert.Contains(t, found, rTest4.ID)
	})

	t.Run("without queue options", func(t *testing.T) {
		rq, err := client.Environments.RunQueue(ctx, orgTest.Name, RunQueueOptions{})
		require.NoError(t, err)

		found := []string{}
		for _, r := range rq.Items {
			found = append(found, r.ID)
		}

		assert.Contains(t, found, rTest1.ID)
		assert.Contains(t, found, rTest2.ID)
		assert.Contains(t, found, rTest3.ID)
		assert.Contains(t, found, rTest4.ID)
		assert.Equal(t, 1, rq.CurrentPage)
		assert.Equal(t, 4, rq.TotalCount)
	})

	t.Run("with queue options", func(t *testing.T) {
		// Request a page number which is out of range. The result should
		// be successful, but return no results if the paging options are
		// properly passed along.
		rq, err := client.Environments.RunQueue(ctx, orgTest.Name, RunQueueOptions{
			ListOptions: ListOptions{
				PageNumber: 999,
				PageSize:   100,
			},
		})
		require.NoError(t, err)

		assert.Empty(t, rq.Items)
		assert.Equal(t, 999, rq.CurrentPage)
		assert.Equal(t, 4, rq.TotalCount)
	})

	t.Run("with invalid name", func(t *testing.T) {
		org, err := client.Environments.Read(ctx, badIdentifier)
		assert.Nil(t, org)
		assert.EqualError(t, err, "invalid value for environment")
	})

	t.Run("when the org does not exist", func(t *testing.T) {
		_, err := client.Environments.Read(ctx, randomString(t))
		assert.Error(t, err)
	})
}
