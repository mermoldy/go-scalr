package scalr

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWorkspaceTagsCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	environment, deleteEnvironment := createEnvironment(t, client)
	defer deleteEnvironment()

	workspace, deleteWorkspace := createWorkspace(t, client, environment)
	defer deleteWorkspace()

	tag, deleteTag := createTag(t, client)
	defer deleteTag()

	t.Run("with valid options", func(t *testing.T) {
		options := WorkspaceTagsCreateOptions{
			WorkspaceID:   workspace.ID,
			WorkspaceTags: []*WorkspaceTag{{ID: tag.ID}},
		}

		err := client.WorkspaceTags.Create(ctx, options)
		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.Workspaces.ReadByID(ctx, workspace.ID)
		require.NoError(t, err)

		for _, item := range refreshed.Tags {
			assert.Equal(t, tag.ID, item.ID)
		}
	})

	t.Run("without valid workspace ID", func(t *testing.T) {
		err := client.WorkspaceTags.Create(ctx, WorkspaceTagsCreateOptions{
			WorkspaceTags: []*WorkspaceTag{{ID: tag.ID}},
		})
		assert.EqualError(t, err, "invalid value for workspace ID")
	})

	t.Run("without valid workspace tags", func(t *testing.T) {
		err := client.WorkspaceTags.Create(ctx, WorkspaceTagsCreateOptions{
			WorkspaceID: workspace.ID,
		})
		assert.EqualError(t, err, "list of tags is required")
	})

	t.Run("when options have an invalid tag", func(t *testing.T) {
		tagID := "tag-invalid-id"
		err := client.WorkspaceTags.Create(ctx, WorkspaceTagsCreateOptions{
			WorkspaceID:   workspace.ID,
			WorkspaceTags: []*WorkspaceTag{{ID: tagID}},
		})
		assert.EqualError(t, err, fmt.Sprintf("Not Found\n\nTag with ID '%s' not found or user unauthorized.", tagID))
	})
}

func TestWorkspaceTagsUpdate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	environment, deleteEnvironment := createEnvironment(t, client)
	defer deleteEnvironment()

	workspace, deleteWorkspace := createWorkspace(t, client, environment)
	defer deleteWorkspace()

	tag, deleteTag := createTag(t, client)
	defer deleteTag()

	t.Run("with valid options", func(t *testing.T) {
		options := WorkspaceTagsUpdateOptions{
			WorkspaceID:   workspace.ID,
			WorkspaceTags: []*WorkspaceTag{{ID: tag.ID}},
		}

		err := client.WorkspaceTags.Update(ctx, options)
		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.Workspaces.ReadByID(ctx, workspace.ID)
		require.NoError(t, err)

		for _, item := range refreshed.Tags {
			assert.Equal(t, tag.ID, item.ID)
		}
	})

	t.Run("without valid workspace ID", func(t *testing.T) {
		err := client.WorkspaceTags.Update(ctx, WorkspaceTagsUpdateOptions{})
		assert.EqualError(t, err, "invalid value for workspace ID")
	})

	t.Run("with invalid workspace tag", func(t *testing.T) {
		tagID := "tag-invalid-id"
		options := WorkspaceTagsUpdateOptions{
			WorkspaceID:   workspace.ID,
			WorkspaceTags: []*WorkspaceTag{{ID: tagID}},
		}

		err := client.WorkspaceTags.Update(ctx, options)
		assert.EqualError(t, err, fmt.Sprintf("Not Found\n\nTag with ID '%s' not found or user unauthorized.", tagID))
	})

	t.Run("when all tags should be removed", func(t *testing.T) {
		err := client.WorkspaceTags.Update(ctx, WorkspaceTagsUpdateOptions{
			WorkspaceID: workspace.ID,
		})
		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.Workspaces.ReadByID(ctx, workspace.ID)
		require.NoError(t, err)
		assert.Empty(t, refreshed.Tags)
	})
}
