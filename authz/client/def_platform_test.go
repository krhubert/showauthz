package client

import (
	"context"
	"testing"

	"rift/assert"
)

func TestPlatform(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tclient, err := StartTestServer(ctx)
	assert.NoError(t, err)

	riftUserId := "alice"
	riftEmail := "alice@rift.com"

	getRiftUserId := "bob"
	getRiftEmail := "bob@getrift.com"

	nonRiftUserId := "charlie"
	nonRiftEmail := "charlie@example.com"

	t.Run("relation_chameleoner", func(t *testing.T) {
		err := tclient.WritePlatfromChameleoner(ctx, riftUserId, riftEmail)
		assert.NoError(t, err)

		err = tclient.WritePlatfromChameleoner(ctx, getRiftUserId, getRiftEmail)
		assert.NoError(t, err)

		err = tclient.WritePlatfromChameleoner(ctx, nonRiftUserId, nonRiftEmail)
		assert.NoError(t, err)
	})

	t.Run("chameleon", func(t *testing.T) {
		err := tclient.CanChameleon(ctx, riftUserId)
		assert.NoError(t, err)

		err = tclient.CanChameleon(ctx, getRiftUserId)
		assert.NoError(t, err)

		err = tclient.CanChameleon(ctx, nonRiftUserId)
		assert.ErrorContains(t, err, &ErrDenied{})
	})
}
