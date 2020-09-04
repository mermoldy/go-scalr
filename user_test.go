package scalr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersReadCurrent(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	u, err := client.Users.ReadCurrent(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, u.ID)
	assert.NotEmpty(t, u.Username)
}
