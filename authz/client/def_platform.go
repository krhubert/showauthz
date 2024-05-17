package client

import (
	"context"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
)

const platformId = "rift"

func (c *Client) WritePlatfromChameleoner(
	ctx context.Context,
	userId string,
	email string,
) error {
	rel := &pb.Relationship{
		Resource: objRef(definitionPlatform, platformId),
		Relation: relationChameleoner,
		Subject:  subRef(definitionUser, userId),
		OptionalCaveat: &pb.ContextualizedCaveat{
			CaveatName: caveatChameleonEmail,
			Context:    newCaveatChameleonEmail(email),
		},
	}

	return c.writeRelationship(ctx, rel)
}

func (c *Client) CanChameleon(ctx context.Context, userId string) error {
	req := &pb.CheckPermissionRequest{
		Resource:    objRef(definitionPlatform, platformId),
		Permission:  permissionChameleon,
		Subject:     subRef(definitionUser, userId),
		Consistency: fullConsistency(),
	}

	return c.checkPermission(ctx, req)
}
