package client

import (
	"context"
	"testing"

	"rift/assert"
)

func TestHoliday(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tclient, err := StartTestServer(ctx)
	assert.NoError(t, err)

	orgId := "rift"
	apiKey := "key"
	adminId := "alice"
	sdrId := "bob"
	holidayId := "holiday"

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

	t.Run("holiday", func(t *testing.T) {
		t.Run("relation_organization", func(t *testing.T) {
			err := tclient.WriteHolidayOrganization(ctx, holidayId, orgId)
			assert.NoError(t, err)
		})

		t.Run("edit", func(t *testing.T) {
			err := tclient.CanEditHoliday(ctx, holidayId, apiKey)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanEditHoliday(ctx, holidayId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanEditHoliday(ctx, holidayId, adminId)
			assert.NoError(t, err)
		})

		t.Run("view", func(t *testing.T) {
			err := tclient.CanViewHoliday(ctx, holidayId, apiKey)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanViewHoliday(ctx, holidayId, sdrId)
			assert.NoError(t, err)

			err = tclient.CanViewHoliday(ctx, holidayId, adminId)
			assert.NoError(t, err)
		})

		t.Run("delete", func(t *testing.T) {
			err := tclient.CanDeleteHoliday(ctx, holidayId, apiKey)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanDeleteHoliday(ctx, holidayId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanDeleteHoliday(ctx, holidayId, adminId)
			assert.NoError(t, err)
		})
	})
}
