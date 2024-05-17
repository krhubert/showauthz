package syncer

import (
	"context"
	"sort"
	"strconv"
	"sync/atomic"
	"testing"

	"rift/assert"
	"rift/authz/client"
	"rift/memdb"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/spicedb/pkg/tuple"
)

type mockMutes struct{}

func (m *mockMutes) Lock(ctx context.Context) error { return nil }
func (m *mockMutes) Unlock(ctx context.Context)     {}

type mockDatabase struct {
	shouldSync    bool
	syncCompleted atomic.Bool
}

func (m *mockDatabase) ShouldSync(ctx context.Context) (bool, error) { return m.shouldSync, nil }
func (m *mockDatabase) SyncCompleted(ctx context.Context) error {
	m.syncCompleted.Store(true)
	return nil
}

func (m *mockDatabase) Members(ctx context.Context) ([]*memdb.Member, error) {
	return []*memdb.Member{
		{ID: "alice", OrganizationID: "rift", Role: memdb.RoleAdmin},
		{ID: "bob", OrganizationID: "rift", Role: memdb.RoleSDR},
	}, nil
}

func (m *mockDatabase) OffDays(ctx context.Context) ([]*memdb.OffDay, error) {
	return []*memdb.OffDay{
		{ID: "offday", OrganizationID: "rift"},
	}, nil
}

type mockLargeDatabase struct {
	mockDatabase
	n int
}

func (m *mockLargeDatabase) Members(ctx context.Context) ([]*memdb.Member, error) {
	members := make([]*memdb.Member, m.n)
	for i := 0; i < m.n; i++ {
		members[i] = &memdb.Member{
			ID:             strconv.Itoa(i),
			OrganizationID: "rift",
			Role:           memdb.RoleAdmin,
		}
	}
	return members, nil
}

func TestSyncerSync(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tclient, err := client.StartTestServer(ctx)
	assert.NoError(t, err)

	md := &mockDatabase{shouldSync: true}
	syncer, err := New(md, &mockMutes{}, tclient)
	assert.NoError(t, err)

	tests := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"Sync", syncer.Sync},
		{"Resync", syncer.Resync},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, tt.fn(ctx))
			assert.True(t, md.syncCompleted.Load())
			md.syncCompleted.Store(false)

			t.Run("organization", func(t *testing.T) {
				rels, err := tclient.ReadRelationships(ctx,
					&pb.RelationshipFilter{ResourceType: "organization"},
				)
				assert.NoError(t, err)

				relStrs := make([]string, len(rels))
				for i, rel := range rels {
					relStrs[i] = tuple.MustStringRelationship(rel)
				}

				sort.Strings(relStrs)
				assert.Equal(t, relStrs, []string{
					`organization:rift#admin@member:alice[products:{"enabled":[]}]`,
					`organization:rift#sdr@member:bob[products:{"enabled":[]}]`,
				})
			})

			t.Run("offday", func(t *testing.T) {
				rels, err := tclient.ReadRelationships(ctx,
					&pb.RelationshipFilter{ResourceType: "offday"},
				)
				assert.NoError(t, err)

				relStrs := make([]string, len(rels))
				for i, rel := range rels {
					relStrs[i] = tuple.MustStringRelationship(rel)
				}

				assert.Equal(t, relStrs, []string{
					`offday:offday#organization@organization:rift`,
				})
			})
		})
	}
}

func TestSyncerShouldSync(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tclient, err := client.StartTestServer(ctx)
	assert.NoError(t, err)

	syncer, err := New(
		&mockDatabase{shouldSync: false},
		&mockMutes{},
		tclient,
	)
	assert.NoError(t, err)

	tests := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"Sync", syncer.Sync},
		{"Resync", syncer.Resync},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, tt.fn(ctx))
			t.Run("organization", func(t *testing.T) {
				rels, err := tclient.ReadRelationships(ctx,
					&pb.RelationshipFilter{ResourceType: "organization"},
				)
				assert.NoError(t, err)
				assert.Len(t, rels, 0)
			})
		})
	}
}

func TestSyncerPerformance(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tclient, err := client.StartTestServer(ctx)
	assert.NoError(t, err)

	syncer, err := New(
		&mockLargeDatabase{
			mockDatabase: mockDatabase{shouldSync: true},
			n:            100_000,
		},
		&mockMutes{},
		tclient,
	)
	assert.NoError(t, err)

	t.Run("sync", func(t *testing.T) {
		assert.NoError(t, syncer.Sync(ctx))
	})
	t.Run("sync again", func(t *testing.T) {
		// check performance of inserting already synced data
		assert.NoError(t, syncer.Sync(ctx))
	})
	t.Run("delete", func(t *testing.T) {
		assert.NoError(t, syncer.delete(ctx))
	})
}
