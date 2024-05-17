package client

import (
	"context"
	"testing"

	"rift/assert"
)

func TestOrganization(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tclient, err := StartTestServer(ctx)
	assert.NoError(t, err)

	orgId := "rift"
	apiKey := "key"
	adminId := "admin"
	sdrId := "bob"

	t.Run("relations", func(t *testing.T) {
		t.Run("apikey", func(t *testing.T) {
			err := tclient.WriteOrganizationApiKey(ctx, orgId, apiKey)
			assert.NoError(t, err)
		})

		t.Run("sdr", func(t *testing.T) {
			err := tclient.WriteOrganizationSDR(ctx, orgId, sdrId)
			assert.NoError(t, err)
		})

		t.Run("admin", func(t *testing.T) {
			err := tclient.WriteOrganizationAdmin(ctx, orgId, adminId)
			assert.NoError(t, err)
		})
	})

	t.Run("permissions", func(t *testing.T) {
		t.Run("edit_settings", func(t *testing.T) {
			err := tclient.CanEditOrganizationSettings(ctx, orgId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanEditOrganizationSettings(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("view_settings", func(t *testing.T) {
			err := tclient.CanViewOrganizationSettings(ctx, orgId, sdrId)
			assert.NoError(t, err)

			err = tclient.CanViewOrganizationSettings(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("invite_member", func(t *testing.T) {
			err := tclient.CanInviteOrganizationMember(ctx, orgId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanInviteOrganizationMember(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("edit_member", func(t *testing.T) {
			err := tclient.CanEditOrganizationMember(ctx, orgId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanEditOrganizationMember(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("delete_member", func(t *testing.T) {
			err := tclient.CanDeleteOrganizationMember(ctx, orgId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanDeleteOrganizationMember(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("create_team", func(t *testing.T) {
			err := tclient.CanCreateOrganizationTeam(ctx, orgId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanCreateOrganizationTeam(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("create_password", func(t *testing.T) {
			err := tclient.CanCreateOrganizationPassword(ctx, orgId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanCreateOrganizationPassword(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("create_offday", func(t *testing.T) {
			err := tclient.CanCreateOrganizationOffDay(ctx, orgId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanCreateOrganizationOffDay(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("create_holiday", func(t *testing.T) {
			err := tclient.CanCreateOrganizationHoliday(ctx, orgId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanCreateOrganizationHoliday(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("create_sequence", func(t *testing.T) {
			err := tclient.CanCreateOrganizationSequence(ctx, orgId, sdrId)
			assert.NoError(t, err)

			err = tclient.CanCreateOrganizationSequence(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("create_inbox", func(t *testing.T) {
			err := tclient.CanCreateOrganizationInbox(ctx, orgId, sdrId)
			assert.ErrorContains(t, err, &ErrDenied{})

			err = tclient.CanCreateOrganizationInbox(ctx, orgId, adminId)
			assert.NoError(t, err)
		})

		t.Run("create_meeting", func(t *testing.T) {
			err := tclient.CanCreateOrganizationMeeting(ctx, orgId, sdrId)
			assert.NoError(t, err)

			err = tclient.CanCreateOrganizationMeeting(ctx, orgId, adminId)
			assert.NoError(t, err)
		})
	})

	t.Run("lookup", func(t *testing.T) {
		org, err := tclient.GetOrganization(ctx, sdrId)
		assert.NoError(t, err)
		assert.Equal(t, org, &Organization{
			Id:             orgId,
			Access:         true,
			CreateHoliday:  false,
			CreateInbox:    false,
			CreateMeeting:  true,
			CreateOffDay:   false,
			CreatePassword: false,
			CreateSequence: true,
			CreateTeam:     false,
			DeleteMember:   false,
			EditMember:     false,
			EditSettings:   false,
			InviteMember:   false,
			ViewSettings:   true,
		})

		org, err = tclient.GetOrganization(ctx, adminId)
		assert.NoError(t, err)
		assert.Equal(t, org, &Organization{
			Id:             orgId,
			Access:         true,
			CreateHoliday:  true,
			CreateInbox:    true,
			CreateMeeting:  true,
			CreateOffDay:   true,
			CreatePassword: true,
			CreateSequence: true,
			CreateTeam:     true,
			DeleteMember:   true,
			EditMember:     true,
			EditSettings:   true,
			InviteMember:   true,
			ViewSettings:   true,
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("apikey", func(t *testing.T) {
			err := tclient.DeleteOrganizationApiKey(ctx, orgId, apiKey)
			assert.NoError(t, err)
		})

		t.Run("sdr", func(t *testing.T) {
			err := tclient.DeleteOrganizationSDR(ctx, orgId, sdrId)
			assert.NoError(t, err)
		})

		t.Run("admin", func(t *testing.T) {
			err := tclient.DeleteOrganizationAdmin(ctx, orgId, adminId)
			assert.NoError(t, err)

			_, err = tclient.GetOrganization(ctx, adminId)
			assert.ErrorContains(t, err, &ErrDenied{})
		})
	})
}
