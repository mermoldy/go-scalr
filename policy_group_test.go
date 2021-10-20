package scalr

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultPolicyGroupID1   = "pgrp-svsu2dqfvtk5qfg"
	defaultPolicyGroupName1 = "Clouds"
	defaultPolicyGroupID2   = "pgrp-svsu1tn68mhevuo"
)

func TestPolicyGroupsList(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("without list options", func(t *testing.T) {
		pgl, err := client.PolicyGroups.List(ctx, PolicyGroupListOptions{})
		require.NoError(t, err)
		pgIDs := make([]string, len(pgl.Items))
		for _, pg := range pgl.Items {
			pgIDs = append(pgIDs, pg.ID)
		}
		assert.Contains(t, pgIDs, defaultPolicyGroupID1)
		assert.Contains(t, pgIDs, defaultPolicyGroupID2)
		assert.Equal(t, 1, pgl.CurrentPage)
		assert.Equal(t, 3, pgl.TotalCount)
	})

	t.Run("with name and account options", func(t *testing.T) {
		pgl, err := client.PolicyGroups.List(ctx, PolicyGroupListOptions{
			Account: defaultAccountID,
			Name:    defaultPolicyGroupName1,
		})
		require.NoError(t, err)
		assert.Equal(t, 1, pgl.CurrentPage)
		assert.Equal(t, 1, pgl.TotalCount)
		assert.Equal(t, defaultPolicyGroupID1, pgl.Items[0].ID)
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
		assert.Equal(t, 3, pgl.TotalCount)
	})

	t.Run("without a valid account", func(t *testing.T) {
		pgl, err := client.PolicyGroups.List(ctx, PolicyGroupListOptions{Account: badIdentifier})
		assert.Len(t, pgl.Items, 0)
		assert.NoError(t, err)
	})
}

func TestPolicyGroupsCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("with empty options", func(t *testing.T) {
		pg, err := client.PolicyGroups.Create(ctx, PolicyGroupCreateOptions{})
		assert.Nil(t, pg)
		assert.EqualError(t, err, "name is required")
	})

	t.Run("without vcs repo options", func(t *testing.T) {
		pg, err := client.PolicyGroups.Create(ctx, PolicyGroupCreateOptions{
			Name:        String("foo"),
			Account:     &Account{ID: defaultAccountID},
			VcsProvider: &VcsProvider{ID: "vcs-123"},
		})
		assert.Nil(t, pg)
		assert.EqualError(t, err, "vcs repo is required")
	})

	t.Run("when options has an invalid account", func(t *testing.T) {
		var accID = "acc-123"
		pg, err := client.PolicyGroups.Create(ctx, PolicyGroupCreateOptions{
			Name:        String("foo"),
			Account:     &Account{ID: accID},
			VcsProvider: &VcsProvider{ID: "vcs-123"},
			VCSRepo: &PolicyGroupVCSRepoOptions{
				Identifier: String("foo/bar"),
			},
		})
		assert.Nil(t, pg)
		assert.EqualError(
			t,
			err,
			ErrResourceNotFound{
				Message: fmt.Sprintf("Clients with ID '%s' not found or user unauthorized", accID),
			}.Error(),
		)
	})

	t.Run("when options has an invalid vcs provider", func(t *testing.T) {
		var vcsID = "vcs-123"
		pg, err := client.PolicyGroups.Create(ctx, PolicyGroupCreateOptions{
			Name:        String("foo"),
			Account:     &Account{ID: defaultAccountID},
			VcsProvider: &VcsProvider{ID: vcsID},
			VCSRepo: &PolicyGroupVCSRepoOptions{
				Identifier: String("foo/bar"),
			},
		})
		assert.Nil(t, pg)
		assert.EqualError(
			t,
			err,
			ErrResourceNotFound{
				Message: fmt.Sprintf("VcsProvider with ID '%s' not found or user unauthorized", vcsID),
			}.Error(),
		)
	})
}

func TestPolicyGroupsRead(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("when the policy group exists", func(t *testing.T) {
		pg, err := client.PolicyGroups.Read(ctx, defaultPolicyGroupID1)
		require.NoError(t, err)
		assert.Equal(t, defaultPolicyGroupID1, pg.ID)

		t.Run("relationships are properly decoded", func(t *testing.T) {
			assert.Equal(t, pg.Account.ID, defaultAccountID)
		})

		t.Run("vcs repo is properly decoded", func(t *testing.T) {
			assert.NotEmpty(t, pg.VCSRepo.Identifier)
			assert.NotEmpty(t, pg.VCSRepo.Branch)
		})

		t.Run("policies are properly decoded", func(t *testing.T) {
			assert.NotEmpty(t, pg.Policies)
		})
	})

	t.Run("when the policy group does not exist", func(t *testing.T) {
		pg, err := client.PolicyGroups.Read(ctx, "pgrp-nonexisting")
		assert.Nil(t, pg)
		assert.Error(t, err)
	})

	t.Run("without a valid policy group ID", func(t *testing.T) {
		pg, err := client.PolicyGroups.Read(ctx, badIdentifier)
		assert.Nil(t, pg)
		assert.EqualError(t, err, "invalid value for policy group ID")
	})
}

func TestPolicyGroupsUpdate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	pgTest, err := client.PolicyGroups.Read(ctx, defaultPolicyGroupID1)
	require.NoError(t, err)

	t.Run("when updating a subset of values", func(t *testing.T) {
		options := PolicyGroupUpdateOptions{
			Name: String("pg-" + randomString(t)),
		}

		pgAfter, err := client.PolicyGroups.Update(ctx, pgTest.ID, options)
		require.NoError(t, err)

		assert.Equal(t, *options.Name, pgAfter.Name)
		assert.Equal(t, pgTest.OpaVersion, pgAfter.OpaVersion)
	})

	t.Run("when an error is returned from the api", func(t *testing.T) {
		pg, err := client.PolicyGroups.Update(ctx, pgTest.ID, PolicyGroupUpdateOptions{
			OpaVersion: String("nonexisting"),
		})
		assert.Nil(t, pg)
		assert.Error(t, err)
	})

	t.Run("when options has an invalid name", func(t *testing.T) {
		pg, err := client.PolicyGroups.Update(ctx, badIdentifier, PolicyGroupUpdateOptions{})
		assert.Nil(t, pg)
		assert.EqualError(t, err, "invalid value for policy group ID")
	})

	// Restore updated policy group
	_, err = client.PolicyGroups.Update(ctx, pgTest.ID, PolicyGroupUpdateOptions{
		Name: String(pgTest.Name),
	},
	)
	require.NoError(t, err)
}

func TestPolicyGroupsDelete(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("when policy group is in use", func(t *testing.T) {
		err := client.PolicyGroups.Delete(ctx, defaultPolicyGroupID1)
		assert.EqualError(
			t, err,
			"Policy group can not be deleted\n\nPolicy group is in use and can not be removed",
		)
	})

	t.Run("without a valid policy group ID", func(t *testing.T) {
		err := client.PolicyGroups.Delete(ctx, badIdentifier)
		assert.EqualError(t, err, "invalid value for policy group ID")
	})
}
