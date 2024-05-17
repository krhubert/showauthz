package client

import (
	"context"
	"testing"

	"rift/assert"
)

func TestPassword(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tclient, err := StartTestServer(ctx)
	assert.NoError(t, err)

	orgId := "rift"
	apiKey := "key"
	adminId := "alice"
	sdrId := "bob"
	passwordId := "password"

	t.Run("organization", func(t *testing.T) {
		t.Run("relation_apikey", func(t *testing.T) {
			err := tclient.WriteOrganizationApiKey(ctx, orgId, apiKey)
			assert.NoError(t, err)
		})

		t.Run("relation_sdr", func(t *testing.T) {
			err := tclient.WriteOrganizationSDR(ctx, orgId, sdrId)
			assert.NoError(t, err)
		})

		t.Run("relation_admin", func(t *testing.T) {
			err := tclient.WriteOrganizationAdmin(ctx, orgId, adminId)
			assert.NoError(t, err)
		})
	})

	t.Run("password", func(t *testing.T) {
		t.Run("relation_organization", func(t *testing.T) {
			err := tclient.WritePasswordOrganization(ctx, passwordId, orgId)
			assert.NoError(t, err)
		})

		t.Run("edit", func(t *testing.T) {
			err := tclient.CanEditPassword(ctx, passwordId, apiKey)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanEditPassword(ctx, passwordId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanEditPassword(ctx, passwordId, adminId)
			assert.NoError(t, err)
		})

		t.Run("view", func(t *testing.T) {
			err := tclient.CanViewPassword(ctx, passwordId, apiKey)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanViewPassword(ctx, passwordId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanViewPassword(ctx, passwordId, adminId)
			assert.NoError(t, err)
		})

		t.Run("delete", func(t *testing.T) {
			err := tclient.CanDeletePassword(ctx, passwordId, apiKey)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanDeletePassword(ctx, passwordId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanDeletePassword(ctx, passwordId, adminId)
			assert.NoError(t, err)
		})
	})
}
