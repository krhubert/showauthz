package client

import (
	"fmt"
	"regexp"
	"sort"
	"testing"

	"rift/assert"
)

const (
	caveatNameExpr     = "([a-z][a-z0-9_]{1,61}[a-z0-9]/)*[a-z][a-z0-9_]{1,62}[a-z0-9]"
	definitionNameExpr = "([a-z][a-z0-9_]{1,61}[a-z0-9]/)*[a-z][a-z0-9_]{1,62}[a-z0-9]"
	relationNameExpr   = "[a-z][a-z0-9_]{1,62}[a-z0-9]"
	permissionNameExpr = "[a-z][a-z0-9_]{1,62}[a-z0-9]"

	caveatExpr     = `caveat\s(?P<name>%s)(.*)\s+{`
	definitionExpr = `definition\s(?P<name>%s)\s{`
	relationExpr   = `\s+relation\s(?P<name>%s):`
	permissionExpr = `\s+permission\s(?P<name>%s)\s=`
)

var (
	caveatRegex     = regexp.MustCompile(fmt.Sprintf(caveatExpr, caveatNameExpr))
	definitionRegex = regexp.MustCompile(fmt.Sprintf(definitionExpr, definitionNameExpr))
	relationRegex   = regexp.MustCompile(fmt.Sprintf(relationExpr, relationNameExpr))
	permissionRegex = regexp.MustCompile(fmt.Sprintf(permissionExpr, permissionNameExpr))
)

func TestSchema(t *testing.T) {
	findUniqueNames := func(regex *regexp.Regexp) []string {
		names := make(map[string]bool)
		for _, match := range regex.FindAllStringSubmatch(schemaV1, -1) {
			names[match[regex.SubexpIndex("name")]] = true
		}
		var unique []string
		for name := range names {
			unique = append(unique, name)
		}

		return unique
	}

	tests := []struct {
		name string
		re   *regexp.Regexp
		want []string
	}{
		{
			"caveat",
			caveatRegex,
			[]string{
				caveatProducts,
				caveatChameleonEmail,
			},
		},
		{
			"definition",
			definitionRegex,
			[]string{
				definitionUser,
				definitionMember,
				definitionTeam,
				definitionPlatform,
				definitionOrganization,
				definitionApiKey,
				definitionContact,
				definitionInbox,
				definitionOffDay,
				definitionHoliday,
				definitionSequence,
				definitionSequenceAction,
				definitionPassword,
				definitionMeeting,
			},
		},
		{
			"relation",
			relationRegex,
			[]string{
				relationAdmin,
				relationApiKey,
				relationAssignee,
				relationChameleoner,
				relationContact,
				relationEditor,
				relationOrganization,
				relationOwner,
				relationSDR,
				relationSender,
				relationSequence,
				relationViewer,
			},
		},
		{
			"permission",
			permissionRegex,
			[]string{
				permissionAccess,
				permissionChameleon,
				permissionCreateCallStep,
				permissionCreateHoliday,
				permissionCreateInbox,
				permissionCreateMeeting,
				permissionCreateOffDay,
				permissionCreatePassword,
				permissionCreateSequence,
				permissionCreateTeam,
				permissionDelete,
				permissionDeleteMember,
				permissionEdit,
				permissionEditMember,
				permissionEditSettings,
				permissionManageSeat,
				permissionInviteMember,
				permissionOrganizationAdmin,
				permissionOrganizationApikey,
				permissionUploadContact,
				permissionView,
				permissionViewSettings,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findUniqueNames(tt.re)
			sort.Strings(got)
			sort.Strings(tt.want)
			assert.Equal(t, tt.want, got)
		})
	}
}
