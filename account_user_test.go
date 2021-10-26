package scalr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountUsersList(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("with empty options", func(t *testing.T) {
		_, err := client.AccountUsers.List(ctx, AccountUserListOptions{})
		require.Error(t, err)
		assert.EqualError(t, err, "either filter[account] or filter[user] is required")
	})

	t.Run("with account option", func(t *testing.T) {
		aul, err := client.AccountUsers.List(ctx, AccountUserListOptions{
			Account: String(defaultAccountID),
		})
		require.NoError(t, err)

		uIDs := make([]string, len(aul.Items))
		for _, au := range aul.Items {
			uIDs = append(uIDs, au.User.ID)
		}
		assert.Equal(t, 1, aul.CurrentPage)
		assert.True(t, aul.TotalCount >= 1)
		assert.Contains(t, uIDs, defaultUserID)
	})

	t.Run("with user option", func(t *testing.T) {
		aul, err := client.AccountUsers.List(ctx, AccountUserListOptions{
			User: String(defaultUserID),
		})
		require.NoError(t, err)

		aIDs := make([]string, len(aul.Items))
		for _, au := range aul.Items {
			aIDs = append(aIDs, au.Account.ID)
		}
		assert.Equal(t, 1, aul.CurrentPage)
		assert.True(t, aul.TotalCount >= 1)
		assert.Contains(t, aIDs, defaultAccountID)
	})

	t.Run("without a valid account", func(t *testing.T) {
		aul, err := client.AccountUsers.List(ctx, AccountUserListOptions{
			Account: String(badIdentifier),
		})
		assert.NoError(t, err)
		assert.Len(t, aul.Items, 0)
	})

	t.Run("without a valid user", func(t *testing.T) {
		aul, err := client.AccountUsers.List(ctx, AccountUserListOptions{
			User: String(badIdentifier),
		})
		assert.NoError(t, err)
		assert.Len(t, aul.Items, 0)
	})
}
