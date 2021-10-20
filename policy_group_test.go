package scalr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicyGroupsList(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	pgTest1, pgTest1Cleanup := createPolicyGroup(t, client)
	defer pgTest1Cleanup()
	pgTest2, pgTest2Cleanup := createPolicyGroup(t, client)
	defer pgTest2Cleanup()

	t.Run("without list options", func(t *testing.T) {
		pgl, err := client.PolicyGroups.List(ctx, PolicyGroupListOptions{})
		require.NoError(t, err)
		pgIDs := make([]string, len(pgl.Items))
		for _, pg := range pgl.Items {
			pgIDs = append(pgIDs, pg.ID)
		}
		assert.Contains(t, pgIDs, pgTest1.ID)
		assert.Contains(t, pgIDs, pgTest2.ID)
		assert.Equal(t, 1, pgl.CurrentPage)
		assert.Equal(t, 2, pgl.TotalCount)
	})

	t.Run("with list options", func(t *testing.T) {
		pgl, err := client.PolicyGroups.List(ctx, PolicyGroupListOptions{
			ListOptions: ListOptions{
				PageNumber: 999,
				PageSize:   100,
			},
		})
		require.NoError(t, err)
		assert.Empty(t, pgl.Items)
		assert.Equal(t, 999, pgl.CurrentPage)
		assert.Equal(t, 2, pgl.TotalCount)
	})

	t.Run("without a valid account", func(t *testing.T) {
		pgl, err := client.PolicyGroups.List(ctx, PolicyGroupListOptions{Account: badIdentifier})
		assert.Len(t, pgl.Items, 0)
		assert.NoError(t, err)
	})
}
