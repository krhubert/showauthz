package client

import (
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/spicedb/pkg/tuple"
	"google.golang.org/protobuf/types/known/structpb"
)

// objRef builds ObjectReference with the given object type and object id.
func objRef(typ string, id string) *pb.ObjectReference {
	return &pb.ObjectReference{
		ObjectType: typ,
		ObjectId:   id,
	}
}

// subRef builds SubjectReference with the given object type and object id.
func subRef(typ, id string, optionalRelations ...string) *pb.SubjectReference {
	optionalRelation := ""
	if len(optionalRelations) > 1 {
		panic("subref: too many optional relations")
	}
	if len(optionalRelations) > 0 {
		optionalRelation = optionalRelations[0]
	}

	return &pb.SubjectReference{
		Object:           objRef(typ, id),
		OptionalRelation: optionalRelation,
	}
}

func fullConsistency() *pb.Consistency {
	return &pb.Consistency{
		Requirement: &pb.Consistency_FullyConsistent{
			FullyConsistent: true,
		},
	}
}

// relstr converts different spicedb structs to string relation.
// It is used to return a clear error message.
func relstr[
	R *pb.CheckPermissionRequest | *pb.LookupResourcesRequest | *pb.Relationship,
](r R) string {
	switch r := (any)(r).(type) {
	case *pb.CheckPermissionRequest:
		return tuple.MustStringRelationship(&pb.Relationship{
			Resource: &pb.ObjectReference{
				ObjectType: r.Resource.ObjectType,
				ObjectId:   r.Resource.ObjectId,
			},
			Relation: r.Permission,
			Subject: &pb.SubjectReference{
				Object: &pb.ObjectReference{
					ObjectType: r.Subject.Object.ObjectType,
					ObjectId:   r.Subject.Object.ObjectId,
				},
			},
		})
	case *pb.LookupResourcesRequest:
		return tuple.MustStringRelationship(&pb.Relationship{
			Resource: &pb.ObjectReference{
				ObjectType: r.ResourceObjectType,
			},
			Relation: r.Permission,
			Subject: &pb.SubjectReference{
				Object: &pb.ObjectReference{
					ObjectType: r.Subject.Object.ObjectType,
					ObjectId:   r.Subject.Object.ObjectId,
				},
			},
		})
	case *pb.Relationship:
		return tuple.MustStringRelationship(r)
	}
	panic("unreachable")
}

// newStructpb converts a map to a structpb.Struct which is used to create a caveat.
func mustNewStructpb(v map[string]any) *structpb.Struct {
	s, err := structpb.NewStruct(v)
	if err != nil {
		panic(err)
	}
	return s
}

func stringsToAny(s []string) []any {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}
