package scalr

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountIPFencingRead(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("account exists", func(t *testing.T) {
		_, err := client.AccountIPAllowLists.Read(ctx, defaultAccountID)
		require.NoError(t, err)
	})

	t.Run("account does not exist", func(t *testing.T) {
		var accId = "notexisting"
		_, err := client.AccountIPAllowLists.Read(ctx, accId)
		assert.Equal(
			t,
			ErrResourceNotFound{
				Message: fmt.Sprintf("Clients with ID '%s' not found or user unauthorized", accId),
			}.Error(),
			err.Error(),
		)
	})

	t.Run("with invalid acc ID", func(t *testing.T) {
		r, err := client.AccountIPAllowLists.Read(ctx, badIdentifier)
		assert.Nil(t, r)
		assert.EqualError(t, err, "invalid value for account ID")
	})
}

func TestAccountIPFencingUpdate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("valid ip allowlist", func(t *testing.T) {
		options := AccountIPAllowlistUpdateOptions{
			IPAllowlist: &[]string{"0.0.0.0/0", "192.168.0.0/24"},
		}
		account, err := client.AccountIPAllowLists.Update(ctx, defaultAccountID, options)
		require.NoError(t, err)
		for i, ip := range account.IPAllowlist {
			assert.Equal(t, ip, (*options.IPAllowlist)[i])
		}
	})

	t.Run("invalid ip allowlist", func(t *testing.T) {
		options := AccountIPAllowlistUpdateOptions{
			IPAllowlist: &[]string{"127.0.00"},
		}
		account, err := client.AccountIPAllowLists.Update(ctx, defaultAccountID, options)
		assert.Nil(t, account)
		assert.EqualError(t, err, "invalid value for ip allowlist entry: 127.0.00")
	})

	t.Run("invalid ip allowlist ipv6", func(t *testing.T) {
		options := AccountIPAllowlistUpdateOptions{
			IPAllowlist: &[]string{"FE80:CD00:0000:0CDE:1257:0000:211E:729C"},
		}
		account, err := client.AccountIPAllowLists.Update(ctx, defaultAccountID, options)
		assert.Nil(t, account)
		assert.EqualError(t, err, "invalid value for ip allowlist entry: FE80:CD00:0000:0CDE:1257:0000:211E:729C")
	})
}
