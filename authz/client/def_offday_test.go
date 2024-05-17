package client

import (
	"context"
	"testing"

	"rift/assert"
)

func TestOffDay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tclient, err := StartTestServer(ctx)
	assert.NoError(t, err)

	orgId := "rift"
	adminId := "alice"
	sdrId := "bob"
	offDayId := "offDay"

	t.Run("relations", func(t *testing.T) {
		t.Run("relation_sdr", func(t *testing.T) {
			err := tclient.WriteOrganizationSDR(ctx, orgId, sdrId)
			assert.NoError(t, err)
		})

		t.Run("relation_admin", func(t *testing.T) {
			err := tclient.WriteOrganizationAdmin(ctx, orgId, adminId)
			assert.NoError(t, err)
		})
		t.Run("relation_organization", func(t *testing.T) {
			err := tclient.WriteOffDayOrganization(ctx, offDayId, orgId)
			assert.NoError(t, err)
		})
	})

	t.Run("permissions", func(t *testing.T) {
		t.Run("edit", func(t *testing.T) {
			err := tclient.CanEditOffDay(ctx, offDayId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			ids, err := tclient.ListEditOffDays(ctx, sdrId)
			assert.NoError(t, err)
			assert.Len(t, ids, 0)

			err = tclient.CanEditOffDay(ctx, offDayId, adminId)
			assert.NoError(t, err)

			ids, err = tclient.ListEditOffDays(ctx, adminId)
			assert.NoError(t, err)
			assert.Equal(t, ids, []string{offDayId})
		})

		t.Run("view", func(t *testing.T) {
			err := tclient.CanViewOffDay(ctx, offDayId, sdrId)
			assert.NoError(t, err)

			ids, err := tclient.ListViewOffDays(ctx, sdrId)
			assert.NoError(t, err)
			assert.Equal(t, ids, []string{offDayId})

			err = tclient.CanViewOffDay(ctx, offDayId, adminId)
			assert.NoError(t, err)

			ids, err = tclient.ListViewOffDays(ctx, adminId)
			assert.NoError(t, err)
			assert.Equal(t, ids, []string{offDayId})
		})

		t.Run("delete", func(t *testing.T) {
			err = tclient.CanDeleteOffDay(ctx, offDayId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanDeleteOffDay(ctx, offDayId, adminId)
			assert.NoError(t, err)
		})
	})

	t.Run("lookup", func(t *testing.T) {
		offDays, err := tclient.ListOffDays(ctx, sdrId)
		assert.NoError(t, err)
		assert.Equal(t, offDays, map[string]*OffDay{
			offDayId: {Id: offDayId, View: true, Edit: false, Delete: false},
		})

		offDays, err = tclient.ListOffDays(ctx, adminId)
		assert.NoError(t, err)
		assert.Equal(t, offDays, map[string]*OffDay{
			offDayId: {Id: offDayId, View: true, Edit: true, Delete: true},
		})
	})

	t.Run("delete", func(t *testing.T) {
		err := tclient.DeleteOffDayOrganization(ctx, offDayId, orgId)
		assert.NoError(t, err)

		offDays, err := tclient.ListOffDays(ctx, adminId)
		assert.NoError(t, err)
		assert.Nil(t, offDays)
	})
}
