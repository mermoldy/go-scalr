package scalr

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEnvironmentTagsAdd(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	environment, deleteEnvironment := createEnvironment(t, client)
	defer deleteEnvironment()

	tag1, deleteTag1 := createTag(t, client)
	defer deleteTag1()
	tag2, deleteTag2 := createTag(t, client)
	defer deleteTag2()
	tag3, deleteTag3 := createTag(t, client)
	defer deleteTag3()

	t.Run("with valid options", func(t *testing.T) {
		err := client.EnvironmentTags.Add(ctx, environment.ID,
			[]*TagRelation{
				{ID: tag1.ID},
				{ID: tag2.ID},
			},
		)
		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.Environments.Read(ctx, environment.ID)
		require.NoError(t, err)
		assert.Len(t, refreshed.Tags, 2)

		tagIDs := make([]string, len(refreshed.Tags))
		for _, tag := range refreshed.Tags {
			tagIDs = append(tagIDs, tag.ID)
		}
		assert.Contains(t, tagIDs, tag1.ID)
		assert.Contains(t, tagIDs, tag2.ID)
	})

	t.Run("add another one", func(t *testing.T) {
		err := client.EnvironmentTags.Add(ctx, environment.ID, []*TagRelation{{ID: tag3.ID}})
		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.Environments.Read(ctx, environment.ID)
		require.NoError(t, err)
		assert.Len(t, refreshed.Tags, 3)

		tagIDs := make([]string, len(refreshed.Tags))
		for _, tag := range refreshed.Tags {
			tagIDs = append(tagIDs, tag.ID)
		}
		assert.Contains(t, tagIDs, tag1.ID)
		assert.Contains(t, tagIDs, tag2.ID)
		assert.Contains(t, tagIDs, tag3.ID)
	})

	t.Run("with invalid tag", func(t *testing.T) {
		tagID := "tag-invalid-id"
		err := client.EnvironmentTags.Add(ctx, environment.ID, []*TagRelation{{ID: tagID}})
		assert.EqualError(t, err, fmt.Sprintf("Not Found\n\nTag with ID '%s' not found or user unauthorized.", tagID))
	})
}

func TestEnvironmentTagsReplace(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	environment, deleteEnvironment := createEnvironment(t, client)
	defer deleteEnvironment()

	tag1, deleteTag1 := createTag(t, client)
	defer deleteTag1()
	tag2, deleteTag2 := createTag(t, client)
	defer deleteTag2()
	tag3, deleteTag3 := createTag(t, client)
	defer deleteTag3()

	assignTagsToEnvironment(t, client, environment, []*Tag{tag1})

	t.Run("with valid options", func(t *testing.T) {
		err := client.EnvironmentTags.Replace(ctx, environment.ID,
			[]*TagRelation{
				{ID: tag2.ID},
				{ID: tag3.ID},
			},
		)
		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.Environments.Read(ctx, environment.ID)
		require.NoError(t, err)
		assert.Len(t, refreshed.Tags, 2)

		tagIDs := make([]string, len(refreshed.Tags))
		for _, tag := range refreshed.Tags {
			tagIDs = append(tagIDs, tag.ID)
		}
		assert.Contains(t, tagIDs, tag2.ID)
		assert.Contains(t, tagIDs, tag3.ID)
	})

	t.Run("with invalid tag", func(t *testing.T) {
		tagID := "tag-invalid-id"
		err := client.EnvironmentTags.Replace(ctx, environment.ID, []*TagRelation{{ID: tagID}})
		assert.EqualError(t, err, fmt.Sprintf("Not Found\n\nTag with ID '%s' not found or user unauthorized.", tagID))
	})

	t.Run("when all tags should be removed", func(t *testing.T) {
		err := client.EnvironmentTags.Replace(ctx, environment.ID, make([]*TagRelation, 0))
		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.Environments.Read(ctx, environment.ID)
		require.NoError(t, err)
		assert.Empty(t, refreshed.Tags)
	})
}

func TestEnvironmentTagsDelete(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	environment, deleteEnvironment := createEnvironment(t, client)
	defer deleteEnvironment()

	tag1, deleteTag1 := createTag(t, client)
	defer deleteTag1()
	tag2, deleteTag2 := createTag(t, client)
	defer deleteTag2()
	tag3, deleteTag3 := createTag(t, client)
	defer deleteTag3()

	assignTagsToEnvironment(t, client, environment, []*Tag{tag1, tag2, tag3})

	t.Run("with valid options", func(t *testing.T) {
		err := client.EnvironmentTags.Delete(ctx, environment.ID,
			[]*TagRelation{
				{ID: tag1.ID},
				{ID: tag2.ID},
			},
		)
		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.Environments.Read(ctx, environment.ID)
		require.NoError(t, err)
		assert.Len(t, refreshed.Tags, 1)
		assert.Equal(t, tag3.ID, refreshed.Tags[0].ID)
	})

	t.Run("with invalid tag", func(t *testing.T) {
		tagID := "tag-invalid-id"
		err := client.EnvironmentTags.Replace(ctx, environment.ID, []*TagRelation{{ID: tagID}})
		assert.EqualError(t, err, fmt.Sprintf("Not Found\n\nTag with ID '%s' not found or user unauthorized.", tagID))
	})
}
