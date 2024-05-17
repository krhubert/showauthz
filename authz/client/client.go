package client

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed schemas/v1.zed
var schemaV1 string

type ErrDenied struct {
	rel string
}

func (e *ErrDenied) Error() string {
	return "access denied: " + e.rel
}

func (e *ErrDenied) Is(target error) bool {
	_, ok := target.(*ErrDenied)
	return ok
}

func IsDenied(err error) bool {
	return errors.Is(err, &ErrDenied{})
}

type Client struct {
	c *authzed.ClientWithExperimental
}

func New(address string, secret string) (*Client, error) {
	client, err := authzed.NewClientWithExperimentalAPIs(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken(secret),
	)
	if err != nil {
		return nil, err
	}

	return &Client{c: client}, nil
}

// UNSAFE_GetClient is a temporary method to get the underlying client
// It is used for testing purposes and by syncer only.
func (c *Client) UNSAFE_GetClient() *authzed.ClientWithExperimental {
	return c.c
}

func (c *Client) MigrateSchema(ctx context.Context) error {
	req := &pb.WriteSchemaRequest{Schema: schemaV1}
	if _, err := c.c.WriteSchema(ctx, req); err != nil {
		return err
	}
	return nil
}

func (c *Client) writeRelationship(ctx context.Context, rel *pb.Relationship) error {
	req := &pb.WriteRelationshipsRequest{
		Updates: []*pb.RelationshipUpdate{
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

func (c *Client) deleteRelationship(ctx context.Context, rel *pb.Relationship) error {
	req := &pb.WriteRelationshipsRequest{
		Updates: []*pb.RelationshipUpdate{
			{
				Operation:    pb.RelationshipUpdate_OPERATION_DELETE,
				Relationship: rel,
			},
		},
	}

	if _, err := c.c.WriteRelationships(ctx, req); err != nil {
		return fmt.Errorf("authz: delete relationships %q: %w", relstr(rel), err)
	}
	return nil
}

func (c *Client) checkPermission(ctx context.Context, req *pb.CheckPermissionRequest) error {
	resp, err := c.c.CheckPermission(ctx, req)
	if err != nil {
		return fmt.Errorf("authz: check permission %q: %w", relstr(req), err)
	}

	if resp.Permissionship != pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
		return &ErrDenied{relstr(req)}
	}

	return nil
}

func (c *Client) lookupResources(ctx context.Context, req *pb.LookupResourcesRequest) ([]string, error) {
	stream, err := c.c.LookupResources(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("authz: lookup resources %q: %w", relstr(req), err)
	}

	var ids []string
	for {
		resp, err := stream.Recv()
		switch {
		case errors.Is(err, io.EOF):
			return ids, nil
		case err != nil:
			return nil, err
		default:
			ids = append(ids, resp.ResourceObjectId)
		}
	}
}

func (c *Client) ReadRelationships(
	ctx context.Context,
	req *pb.RelationshipFilter,
) ([]*pb.Relationship, error) {
	stream, err := c.c.ReadRelationships(ctx, &pb.ReadRelationshipsRequest{
		Consistency:        fullConsistency(),
		RelationshipFilter: req,
	})
	if err != nil {
		return nil, fmt.Errorf("authz: read relationships: %w", err)
	}

	var rels []*pb.Relationship
	for {
		resp, err := stream.Recv()
		switch {
		case errors.Is(err, io.EOF):
			return rels, nil
		case err != nil:
			return nil, err
		default:
			rels = append(rels, resp.Relationship)
		}
	}
}
