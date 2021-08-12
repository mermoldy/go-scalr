package scalr

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVCSRevisionRead(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("when the vcs revision does not exist", func(t *testing.T) {
		var vcsId = "nonexisting"
		cv, err := client.VcsRevisions.Read(ctx, vcsId)
		assert.Nil(t, cv)
		assert.Equal(
			t,
			err,
			errors.New(
				fmt.Sprintf("VcsRevisionBinding with ID '%s' not found or user unauthorized", vcsId),
			),
		)
	})

	t.Run("with invalid vcs revision id", func(t *testing.T) {
		cv, err := client.VcsRevisions.Read(ctx, badIdentifier)
		assert.Nil(t, cv)
		assert.EqualError(t, err, "invalid value for vcs revision ID")
	})
}
