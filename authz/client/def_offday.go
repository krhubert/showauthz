package client

import (
	"context"
	"fmt"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
)

type OffDay struct {
	Id     string
	View   bool
	Edit   bool
	Delete bool
}

func RelationOffDayOrganization(
	offDayId string,
	organizationId string,
) *pb.Relationship {
	return &pb.Relationship{
		Resource: objRef(definitionOffDay, offDayId),
		Relation: relationOrganization,
		Subject:  subRef(definitionOrganization, organizationId),
	}
}

func (c *Client) WriteOffDayOrganization(
	ctx context.Context,
	offDayId string,
	organizationId string,
) error {
	rel := RelationOffDayOrganization(offDayId, organizationId)
	return c.writeRelationship(ctx, rel)
}

func (c *Client) DeleteOffDayOrganization(
	ctx context.Context,
	offDayId string,
	organizationId string,
) error {
	rel := RelationOffDayOrganization(offDayId, organizationId)
	return c.deleteRelationship(ctx, rel)
}

func (c *Client) CanEditOffDay(
	ctx context.Context,
	offDayId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOffDay, offDayId),
		Permission:  permissionEdit,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}

func (c *Client) ListEditOffDays(
	ctx context.Context,
	memberId string,
) ([]string, error) {
	req := &pb.LookupResourcesRequest{
		ResourceObjectType: definitionOffDay,
		Permission:         permissionEdit,
		Subject:            subRef(definitionMember, memberId),
		Context:            newCaveatProductsAllow(),
		Consistency:        fullConsistency(),
	}
	return c.lookupResources(ctx, req)
}

func (c *Client) CanViewOffDay(
	ctx context.Context,
	offDayId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOffDay, offDayId),
		Permission:  permissionView,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}

func (c *Client) ListViewOffDays(
	ctx context.Context,
	memberId string,
) ([]string, error) {
	req := &pb.LookupResourcesRequest{
		ResourceObjectType: definitionOffDay,
		Permission:         permissionView,
		Subject:            subRef(definitionMember, memberId),
		Context:            newCaveatProductsAllow(),
		Consistency:        fullConsistency(),
	}
	return c.lookupResources(ctx, req)
}

func (c *Client) CanDeleteOffDay(
	ctx context.Context,
	offDayId string,
	memberId string,
) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionOffDay, offDayId),
		Permission:  permissionDelete,
		Subject:     subRef(definitionMember, memberId),
		Consistency: fullConsistency(),
		Context:     newCaveatProductsAllow(),
	}
	return c.checkPermission(ctx, req)
}

func (c *Client) ListDeleteOffDays(
	ctx context.Context,
	memberId string,
) ([]string, error) {
	req := &pb.LookupResourcesRequest{
		ResourceObjectType: definitionOffDay,
		Permission:         permissionDelete,
		Subject:            subRef(definitionMember, memberId),
		Context:            newCaveatProductsAllow(),
		Consistency:        fullConsistency(),
	}
	return c.lookupResources(ctx, req)
}

func (c *Client) ListOffDays(ctx context.Context, memberId string) (map[string]*OffDay, error) {
	viewIds, err := c.ListViewOffDays(ctx, memberId)
	if err != nil {
		return nil, err
	}

	if len(viewIds) == 0 {
		return nil, nil
	}

	offDays := map[string]*OffDay{}
	items := make([]*pb.BulkCheckPermissionRequestItem, 0, 2*len(viewIds))
	for _, viewId := range viewIds {
		offDays[viewId] = &OffDay{
			Id:   viewId,
			View: true,
		}

		items = append(
			items,
			&pb.BulkCheckPermissionRequestItem{
				Resource:   objRef(definitionOffDay, viewId),
				Permission: permissionEdit,
				Subject:    subRef(definitionMember, memberId),
				Context:    newCaveatProductsAllow(),
			},
			&pb.BulkCheckPermissionRequestItem{
				Resource:   objRef(definitionOffDay, viewId),
				Permission: permissionDelete,
				Subject:    subRef(definitionMember, memberId),
				Context:    newCaveatProductsAllow(),
			},
		)
	}
	req := &pb.BulkCheckPermissionRequest{
		Consistency: fullConsistency(),
		Items:       items,
	}
	resp, err := c.c.BulkCheckPermission(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, pair := range resp.GetPairs() {
		switch r := pair.GetResponse().(type) {
		case *pb.BulkCheckPermissionPair_Item:
			if r.Item.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
				switch pair.Request.Permission {
				case permissionEdit:
					offDays[pair.Request.Resource.ObjectId].Edit = true
				case permissionDelete:
					offDays[pair.Request.Resource.ObjectId].Delete = true
				default:
					return nil, fmt.Errorf("unexpected permission %s", pair.Request.Permission)
				}
			}
		case *pb.BulkCheckPermissionPair_Error:
			return nil, fmt.Errorf("%d %s", r.Error.Code, r.Error.Message)
		default:
			return nil, fmt.Errorf("unexpected response type %T", r)
		}
	}

	return offDays, nil
}
