package scalr

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVCSProvidersList(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	vcsTest1, vcsTest1Cleanup := createVcsProvider(t, client, nil)
	defer vcsTest1Cleanup()

	env, envCleanup := createEnvironment(t, client)
	defer envCleanup()
	vcsTest2, vcsTest2Cleanup := createVcsProvider(t, client, env)
	defer vcsTest2Cleanup()

	t.Run("without list options", func(t *testing.T) {
		response, err := client.VcsProviders.List(ctx, VcsProvidersListOptions{})
		require.NoError(t, err)
		vcsIDs := make([]string, len(response.Items))
		for _, vcs := range response.Items {
			vcsIDs = append(vcsIDs, vcs.ID)
		}
		assert.Contains(t, vcsIDs, vcsTest1.ID)
		assert.Contains(t, vcsIDs, vcsTest2.ID)
		assert.Equal(t, 1, response.CurrentPage)
	})

	t.Run("with list options", func(t *testing.T) {
		response, err := client.VcsProviders.List(
			ctx,
			VcsProvidersListOptions{Query: &vcsTest2.Name, VcsType: &vcsTest2.VcsType, Environment: &env.ID},
		)
		require.NoError(t, err)
		vcsIDs := make([]string, len(response.Items))
		for _, vcs := range response.Items {
			vcsIDs = append(vcsIDs, vcs.ID)
		}
		assert.Contains(t, vcsIDs, vcsTest2.ID)
		assert.Equal(t, 1, response.CurrentPage)
		assert.Equal(t, 1, response.TotalCount)
	})

	t.Run("with invalid environment filter", func(t *testing.T) {
		response, err := client.VcsProviders.List(ctx, VcsProvidersListOptions{Environment: String(badIdentifier)})
		assert.Len(t, response.Items, 0)
		assert.NoError(t, err)
	})
}
