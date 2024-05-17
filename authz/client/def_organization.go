package client

import (
	"context"
	"fmt"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
)

type Organization struct {
	Id             string
	Access         bool
	EditSettings   bool
	ViewSettings   bool
	InviteMember   bool
	EditMember     bool
	DeleteMember   bool
	CreateTeam     bool
	CreatePassword bool
	CreateOffDay   bool
	CreateHoliday  bool
	CreateSequence bool
	CreateInbox    bool
	CreateMeeting  bool
}

func RelationOrganizationApiKey(
	organizationId string,
	apiKeyId string,
) *pb.Relationship {
	return &pb.Relationship{
		Resource: objRef(definitionOrganization, organizationId),
		Relation: relationApiKey,
		Subject:  subRef(definitionApiKey, apiKeyId),
	}
}

func (c *Client) WriteOrganizationApiKey(
	ctx context.Context,
	organizationId string,
	apiKeyId string,
) error {
	rel := RelationOrganizationApiKey(organizationId, apiKeyId)
	return c.writeRelationship(ctx, rel)
}

func (c *Client) DeleteOrganizationApiKey(
	ctx context.Context,
	organizationId string,
	apiKeyId string,
) error {
	rel := RelationOrganizationApiKey(organizationId, apiKeyId)
	return c.deleteRelationship(ctx, rel)
}

func RelationOrganizationAdmin(
	organizationId string,
	memberId string,
	products ...string,
) *pb.Relationship {
	return &pb.Relationship{
		Resource: objRef(definitionOrganization, organizationId),
		Relation: relationAdmin,
		Subject:  subRef(definitionMember, memberId),
		OptionalCaveat: &pb.ContextualizedCaveat{
			CaveatName: caveatProducts,
			Context:    newCaveatProducts(products),
		},
	}
}

func (c *Client) WriteOrganizationAdmin(
	ctx context.Context,
	organizationId string,
	memberId string,
	products ...string,
) error {
	// NOTE: write organization admin and sdr are different
	// than other relationships, because they are mutually exclusive.
	// this is why there's a deletion of all other roles before
	// writing the new one.
	rel := RelationOrganizationAdmin(organizationId, memberId)
	req := &pb.WriteRelationshipsRequest{
		Updates: []*pb.RelationshipUpdate{
			{
				Operation:    pb.RelationshipUpdate_OPERATION_DELETE,
				Relationship: RelationOrganizationSDR(organizationId, memberId),
			},
			{
				Operation:    pb.RelationshipUpdate_OPERATION_TOUCH,
				Relationship: rel,
			},
		},
	}

	if _, err := c.c.WriteRelationships(ctx, req); err != nil {
		return fmt.Errorf("authz: write relationships %q: %w", relstr(rel), err)
	}
	return nil
}

func (c *Client) DeleteOrganizationAdmin(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	rel := RelationOrganizationAdmin(organizationId, memberId)
	return c.deleteRelationship(ctx, rel)
}

func RelationOrganizationSDR(
	organizationId string,
	memberId string,
	products ...string,
) *pb.Relationship {
	return &pb.Relationship{
		Resource: objRef(definitionOrganization, organizationId),
		Relation: relationSDR,
		Subject:  subRef(definitionMember, memberId),
		OptionalCaveat: &pb.ContextualizedCaveat{
			CaveatName: caveatProducts,
			Context:    newCaveatProducts(products),
		},
	}
}

func (c *Client) WriteOrganizationSDR(
	ctx context.Context,
	organizationId string,
	memberId string,
	products ...string,
) error {
	rel := RelationOrganizationSDR(organizationId, memberId)
	req := &pb.WriteRelationshipsRequest{
		Updates: []*pb.RelationshipUpdate{
			{
				Operation:    pb.RelationshipUpdate_OPERATION_DELETE,
				Relationship: RelationOrganizationAdmin(organizationId, memberId),
			},
			{
				Operation:    pb.RelationshipUpdate_OPERATION_TOUCH,
				Relationship: rel,
			},
		},
	}
	if _, err := c.c.WriteRelationships(ctx, req); err != nil {
		return fmt.Errorf("authz: write relationships %q: %w", relstr(rel), err)
	}
	return nil
}

func (c *Client) DeleteOrganizationSDR(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	rel := RelationOrganizationSDR(organizationId, memberId)
	return c.deleteRelationship(ctx, rel)
}

func (c *Client) CanEditOrganizationSettings(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionEditSettings,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) CanViewOrganizationSettings(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionViewSettings,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) CanInviteOrganizationMember(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionInviteMember,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) CanEditOrganizationMember(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionEditMember,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) CanDeleteOrganizationMember(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionDeleteMember,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

// create_team
func (c *Client) CanCreateOrganizationTeam(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionCreateTeam,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) CanCreateOrganizationPassword(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionCreatePassword,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) CanCreateOrganizationOffDay(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionCreateOffDay,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) CanCreateOrganizationHoliday(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionCreateHoliday,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) CanCreateOrganizationSequence(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionCreateSequence,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) CanCreateOrganizationInbox(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionCreateInbox,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) CanCreateOrganizationMeeting(
	ctx context.Context,
	organizationId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOrganization, organizationId),
		Permission:  permissionCreateMeeting,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}

	return c.checkPermission(ctx, req)
}

func (c *Client) GetOrganization(ctx context.Context, memberId string) (*Organization, error) {
	req := &pb.LookupResourcesRequest{
		ResourceObjectType: definitionOrganization,
		Permission:         permissionAccess,
		Subject:            subRef(definitionMember, memberId),
		Context:            newCaveatProductsAllow(),
		Consistency:        fullConsistency(),
	}

	orgIds, err := c.lookupResources(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(orgIds) == 0 {
		return nil, &ErrDenied{relstr(req)}
	}
	if len(orgIds) > 1 {
		return nil, fmt.Errorf("authz: internal error %q: member belons to many organizations", relstr(req))
	}

	checkItem := func(permission string) *pb.BulkCheckPermissionRequestItem {
		return &pb.BulkCheckPermissionRequestItem{
			Resource:   objRef(definitionOrganization, orgIds[0]),
			Permission: permission,
			Subject:    subRef(definitionMember, memberId),
			Context:    newCaveatProductsAllow(),
		}
	}

	items := []*pb.BulkCheckPermissionRequestItem{
		checkItem(permissionEditSettings),
		checkItem(permissionViewSettings),
		checkItem(permissionInviteMember),
		checkItem(permissionEditMember),
		checkItem(permissionDeleteMember),
		checkItem(permissionCreateTeam),
		checkItem(permissionCreatePassword),
		checkItem(permissionCreateOffDay),
		checkItem(permissionCreateHoliday),
		checkItem(permissionCreateSequence),
		checkItem(permissionCreateInbox),
		checkItem(permissionCreateMeeting),
	}

	reqC := &pb.BulkCheckPermissionRequest{
		Consistency: fullConsistency(),
		Items:       items,
	}
	resp, err := c.c.BulkCheckPermission(ctx, reqC)
	if err != nil {
		return nil, err
	}

	org := Organization{
		Id:     orgIds[0],
		Access: true,
	}

	for _, pair := range resp.GetPairs() {
		switch r := pair.GetResponse().(type) {
		case *pb.BulkCheckPermissionPair_Item:
			if r.Item.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
				switch pair.Request.Permission {
				case permissionEditSettings:
					org.EditSettings = true
				case permissionViewSettings:
					org.ViewSettings = true
				case permissionInviteMember:
					org.InviteMember = true
				case permissionEditMember:
					org.EditMember = true
				case permissionDeleteMember:
					org.DeleteMember = true
				case permissionCreateTeam:
					org.CreateTeam = true
				case permissionCreatePassword:
					org.CreatePassword = true
				case permissionCreateOffDay:
					org.CreateOffDay = true
				case permissionCreateHoliday:
					org.CreateHoliday = true
				case permissionCreateSequence:
					org.CreateSequence = true
				case permissionCreateInbox:
					org.CreateInbox = true
				case permissionCreateMeeting:
					org.CreateMeeting = true
				default:
					return nil, fmt.Errorf("unexpected permission %q", pair.Request.Permission)
				}
			}
		case *pb.BulkCheckPermissionPair_Error:
			return nil, fmt.Errorf("%d %s", r.Error.Code, r.Error.Message)
		default:
			return nil, fmt.Errorf("unexpected response type %T", r)
		}
	}

	return &org, nil
}
