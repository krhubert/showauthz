package client

import (
	"context"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
)

func (c *Client) WriteTeamOrganization(
	ctx context.Context,
	teamId string,
	organizationId string,
) error {
	rel := &pb.Relationship{
		Resource: objRef(definitionTeam, teamId),
		Relation: relationOrganization,
		Subject:  subRef(definitionOrganization, organizationId),
	}

	return c.writeRelationship(ctx, rel)
}

func (c *Client) CanEditTeam(
	ctx context.Context,
	teamId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionTeam, teamId),
		Permission:  permissionEdit,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}

func (c *Client) CanViewTeam(
	ctx context.Context,
	teamId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionTeam, teamId),
		Permission:  permissionView,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}

func (c *Client) CanDeleteTeam(
	ctx context.Context,
	teamId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionTeam, teamId),
		Permission:  permissionDelete,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}
