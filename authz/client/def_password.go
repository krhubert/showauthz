package client

import (
	"context"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
)

func (c *Client) WritePasswordOrganization(
	ctx context.Context,
	passwordId string,
	organizationId string,
) error {
	rel := &pb.Relationship{
		Resource: objRef(definitionPassword, passwordId),
		Relation: relationOrganization,
		Subject:  subRef(definitionOrganization, organizationId),
	}

	return c.writeRelationship(ctx, rel)
}

func (c *Client) CanEditPassword(
	ctx context.Context,
	passwordId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionPassword, passwordId),
		Permission:  permissionEdit,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}

func (c *Client) CanViewPassword(
	ctx context.Context,
	passwordId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionPassword, passwordId),
		Permission:  permissionView,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}

func (c *Client) CanDeletePassword(
	ctx context.Context,
	passwordId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionPassword, passwordId),
		Permission:  permissionDelete,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}
