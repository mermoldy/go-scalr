package scalr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVCSRevisionRead(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("when the vcs revision does not exist", func(t *testing.T) {
		cv, err := client.VcsRevisions.Read(ctx, "nonexisting")
		assert.Nil(t, cv)
		assert.Equal(t, err, ErrResourceNotFound)
	})

	t.Run("with invalid vcs revision id", func(t *testing.T) {
		cv, err := client.VcsRevisions.Read(ctx, badIdentifier)
		assert.Nil(t, cv)
		assert.EqualError(t, err, "invalid value for vcs revision ID")
	})
}
