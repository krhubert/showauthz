package client

import (
	"context"
	"testing"

	"rift/assert"
)

func TestTeam(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tclient, err := StartTestServer(ctx)
	assert.NoError(t, err)

	orgId := "rift"
	apiKey := "key"
	adminId := "alice"
	sdrId := "bob"
	teamId := "team"

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

	t.Run("team", func(t *testing.T) {
		t.Run("relation_organization", func(t *testing.T) {
			err := tclient.WriteTeamOrganization(ctx, teamId, orgId)
			assert.NoError(t, err)
		})

		t.Run("edit", func(t *testing.T) {
			err := tclient.CanEditTeam(ctx, teamId, apiKey)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanEditTeam(ctx, teamId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanEditTeam(ctx, teamId, adminId)
			assert.NoError(t, err)
		})

		t.Run("view", func(t *testing.T) {
			err := tclient.CanViewTeam(ctx, teamId, apiKey)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanViewTeam(ctx, teamId, sdrId)
			assert.NoError(t, err)

			err = tclient.CanViewTeam(ctx, teamId, adminId)
			assert.NoError(t, err)
		})

		t.Run("delete", func(t *testing.T) {
			err := tclient.CanDeleteTeam(ctx, teamId, apiKey)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanDeleteTeam(ctx, teamId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanDeleteTeam(ctx, teamId, adminId)
			assert.NoError(t, err)
		})
	})
}
