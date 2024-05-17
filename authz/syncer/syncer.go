package syncer

import (
	"context"
	"fmt"

	"rift/authz/client"
	"rift/memdb"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
)

const batchSize = 1000

type Syncer struct {
	db     Database
	mx     Mutex
	authzC *authzed.ClientWithExperimental
}

func New(db Database, mx Mutex, c *client.Client) (*Syncer, error) {
	return &Syncer{
		db:     db,
		mx:     mx,
		authzC: c.UNSAFE_GetClient(),
	}, nil
}

// Sync syncs the database with the authzed server.
func (s *Syncer) Sync(ctx context.Context) error {
	if err := s.mx.Lock(ctx); err != nil {
		// already locked by other instance
		if isMutexLocked(err) {
			return nil
		}
		return fmt.Errorf("failed to lock: %w", err)
	}
	defer s.mx.Unlock(ctx)

	should, err := s.db.ShouldSync(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if should sync: %w", err)
	}

	if !should {
		return nil
	}

	if err := s.sync(ctx); err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	if err := s.db.SyncCompleted(ctx); err != nil {
		return fmt.Errorf("failed to mark sync completed: %w", err)
	}

	return nil
}

// Resync first removes all existing relationships and then syncs the data again.
// NOTE: deleting all relationships is very expensive and should be used with caution.
// This operation can last up to hours when run on a large dataset - millions of relationships.
func (s *Syncer) Resync(ctx context.Context) error {
	if err := s.mx.Lock(ctx); err != nil {
		// already locked by other instance
		if isMutexLocked(err) {
			return nil
		}
		return fmt.Errorf("failed to lock: %w", err)
	}
	defer s.mx.Unlock(ctx)

	should, err := s.db.ShouldSync(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if should sync: %w", err)
	}

	if !should {
		return nil
	}

	if err := s.delete(ctx); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	if err := s.sync(ctx); err != nil {
		return fmt.Errorf("failed to re-sync: %w", err)
	}

	if err := s.db.SyncCompleted(ctx); err != nil {
		return fmt.Errorf("failed to mark re-sync completed: %w", err)
	}
	return nil
}

func (s *Syncer) delete(ctx context.Context) error {
	deleteReq := func(resourceType string) *pb.DeleteRelationshipsRequest {
		return &pb.DeleteRelationshipsRequest{
			RelationshipFilter: &pb.RelationshipFilter{
				ResourceType: resourceType,
			},
		}
	}

	if _, err := s.authzC.DeleteRelationships(ctx, deleteReq("organization")); err != nil {
		return fmt.Errorf("failed to delete relationships: %w", err)
	}

	if _, err := s.authzC.DeleteRelationships(ctx, deleteReq("offday")); err != nil {
		return fmt.Errorf("failed to delete relationships: %w", err)
	}

	return nil
}

func (s *Syncer) sync(ctx context.Context) error {
	if err := s.syncMembers(ctx); err != nil {
		return fmt.Errorf("failed to sync members: %w", err)
	}

	if err := s.syncOffDays(ctx); err != nil {
		return fmt.Errorf("failed to sync off days: %w", err)
	}

	return nil
}

func (s *Syncer) syncMembers(ctx context.Context) error {
	members, err := s.db.Members(ctx)
	if err != nil {
		return fmt.Errorf("failed to query members: %w", err)
	}

	for i := 0; i < len(members); i += batchSize {
		end := i + batchSize
		if end > len(members) {
			end = len(members)
		}

		req := &pb.WriteRelationshipsRequest{
			Updates: make([]*pb.RelationshipUpdate, 0, batchSize),
		}

		for _, m := range members[i:end] {
			switch {
			case m.Role == memdb.RoleAdmin:
				req.Updates = append(req.Updates, &pb.RelationshipUpdate{
					Operation:    pb.RelationshipUpdate_OPERATION_TOUCH,
					Relationship: client.RelationOrganizationAdmin(m.OrganizationID, m.ID),
				})
			case m.Role == memdb.RoleSDR:
				req.Updates = append(req.Updates, &pb.RelationshipUpdate{
					Operation:    pb.RelationshipUpdate_OPERATION_TOUCH,
					Relationship: client.RelationOrganizationSDR(m.OrganizationID, m.ID),
				})
			default:
				return fmt.Errorf("unknown role: %s", m.Role)
			}
		}

		if _, err := s.authzC.WriteRelationships(ctx, req); err != nil {
			return fmt.Errorf("failed to write relationships: %w", err)
		}
	}

	return nil
}

func (s *Syncer) syncOffDays(ctx context.Context) error {
	offDays, err := s.db.OffDays(ctx)
	if err != nil {
		return fmt.Errorf("failed to query off days: %w", err)
	}

	for i := 0; i < len(offDays); i += batchSize {
		end := i + batchSize
		if end > len(offDays) {
			end = len(offDays)
		}
		req := &pb.WriteRelationshipsRequest{
			Updates: make([]*pb.RelationshipUpdate, 0, batchSize),
		}
		for _, d := range offDays[i:end] {
			req.Updates = append(req.Updates, &pb.RelationshipUpdate{
				Operation:    pb.RelationshipUpdate_OPERATION_TOUCH,
				Relationship: client.RelationOffDayOrganization(d.ID, d.OrganizationID),
			})
		}

		if _, err := s.authzC.WriteRelationships(ctx, req); err != nil {
			return fmt.Errorf("failed to write relationships: %w", err)
		}
	}

	return nil
}
