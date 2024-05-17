// This file contains all definitions, relations,
// permissions, and caveats used in schema.
// Together with the schema_test it guarantees that
// there are no unused definitions, relations, permissions, and caveats.
// Because all of them are not exported,
// the unused linter will catch any unused values.
package client

import (
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	caveatProducts            = "products"
	caveatProductsEnabledArg  = "enabled"  // list<string>
	caveatProductsRequiredArg = "required" // list<string>
	caveatProductsOneOfArg    = "oneOf"    // bool

	caveatChameleonEmail    = "chameleon_email"
	caveatChameleonEmailArg = "email" // string
)

func newCaveatProducts(products []string) *structpb.Struct {
	return mustNewStructpb(map[string]any{
		caveatProductsEnabledArg: stringsToAny(products),
	})
}

func newCaveatProductsAllow() *structpb.Struct {
	return newCaveatProductsCheck([]string{}, false)
}

func newCaveatProductsCheck(required []string, oneOf bool) *structpb.Struct {
	return mustNewStructpb(map[string]interface{}{
		caveatProductsRequiredArg: stringsToAny(required),
		caveatProductsOneOfArg:    oneOf,
	})
}

func newCaveatChameleonEmail(email string) *structpb.Struct {
	return mustNewStructpb(map[string]interface{}{
		caveatChameleonEmailArg: email,
	})
}

const (
	relationAdmin        = "admin"
	relationApiKey       = "apikey"
	relationAssignee     = "assignee"
	relationChameleoner  = "chameleoner"
	relationContact      = "contact"
	relationEditor       = "editor"
	relationOrganization = "organization"
	relationOwner        = "owner"
	relationSDR          = "sdr"
	relationSender       = "sender"
	relationSequence     = "sequence"
	relationViewer       = "viewer"
)

const (
	permissionAccess             = "access"
	permissionChameleon          = "chameleon"
	permissionCreateCallStep     = "create_call_step"
	permissionCreateHoliday      = "create_holiday"
	permissionCreateInbox        = "create_inbox"
	permissionCreateMeeting      = "create_meeting"
	permissionCreateOffDay       = "create_offday"
	permissionCreatePassword     = "create_password"
	permissionCreateSequence     = "create_sequence"
	permissionCreateTeam         = "create_team"
	permissionDelete             = "delete"
	permissionDeleteMember       = "delete_member"
	permissionEdit               = "edit"
	permissionEditMember         = "edit_member"
	permissionEditSettings       = "edit_settings"
	permissionManageSeat         = "manage_seat"
	permissionInviteMember       = "invite_member"
	permissionOrganizationAdmin  = "organization_admin"
	permissionOrganizationApikey = "organization_apikey"
	permissionUploadContact      = "upload_contact"
	permissionView               = "view"
	permissionViewSettings       = "view_settings"
)

const (
	definitionUser           = "user"
	definitionMember         = "member"
	definitionTeam           = "team"
	definitionPlatform       = "platform"
	definitionOrganization   = "organization"
	definitionApiKey         = "apikey"
	definitionContact        = "contact"
	definitionInbox          = "inbox"
	definitionOffDay         = "offday"
	definitionHoliday        = "holiday"
	definitionSequence       = "sequence"
	definitionSequenceAction = "sequence/action"
	definitionPassword       = "password"
	definitionMeeting        = "meeting"
)
