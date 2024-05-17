package client

import (
	"context"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
)

func (c *Client) WriteHolidayOrganization(
	ctx context.Context,
	holidayId string,
	organizationId string,
) error {
	rel := &pb.Relationship{
		Resource: objRef(definitionHoliday, holidayId),
		Relation: relationOrganization,
		Subject:  subRef(definitionOrganization, organizationId),
	}

	return c.writeRelationship(ctx, rel)
}

func (c *Client) CanEditHoliday(
	ctx context.Context,
	holidayId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionHoliday, holidayId),
		Permission:  permissionEdit,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}

func (c *Client) CanViewHoliday(
	ctx context.Context,
	holidayId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionHoliday, holidayId),
		Permission:  permissionView,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}

func (c *Client) CanDeleteHoliday(
	ctx context.Context,
	holidayId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionHoliday, holidayId),
		Permission:  permissionDelete,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}
