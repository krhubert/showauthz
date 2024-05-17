package memdb

import "sync"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleSDR   Role = "sdr"
)

type Member struct {
	ID             string
	OrganizationID string
	Role           Role
}

type OffDay struct {
	ID             string
	OrganizationID string
}

type DB struct {
	members map[string]*Member
	offDays map[string]*OffDay

	mx sync.RWMutex
}

func New() *DB {
	return &DB{
		members: make(map[string]*Member),
		offDays: make(map[string]*OffDay),
	}
}

func (db *DB) AddMember(m *Member) {
	db.mx.Lock()
	db.members[m.ID] = m
	db.mx.Unlock()
}

func (db *DB) GetMember(id string) *Member {
	db.mx.RLock()
	m := db.members[id]
	db.mx.RUnlock()
	return m
}

func (db *DB) AllMembers() []*Member {
	db.mx.RLock()
	members := make([]*Member, 0, len(db.members))
	for _, m := range db.members {
		members = append(members, m)
	}
	db.mx.RUnlock()
	return members
}

func (db *DB) GetMembers(ids ...string) []*Member {
	db.mx.RLock()
	members := make([]*Member, 0, len(ids))
	for _, id := range ids {
		members = append(members, db.members[id])
	}
	db.mx.RUnlock()
	return members
}

func (db *DB) DeleteMember(id string) {
	db.mx.Lock()
	delete(db.members, id)
	db.mx.Unlock()
}

func (db *DB) AddOffDay(o *OffDay) {
	db.mx.Lock()
	db.offDays[o.ID] = o
	db.mx.Unlock()
}

func (db *DB) GetOffDay(id string) *OffDay {
	db.mx.RLock()
	o := db.offDays[id]
	db.mx.RUnlock()
	return o
}

func (db *DB) AllOffDays() []*OffDay {
	db.mx.RLock()
	offDays := make([]*OffDay, 0, len(db.offDays))
	for _, o := range db.offDays {
		offDays = append(offDays, o)
	}
	db.mx.RUnlock()
	return offDays
}

func (db *DB) GetOffDays(ids ...string) []*OffDay {
	db.mx.RLock()
	offDays := make([]*OffDay, 0, len(ids))
	for _, id := range ids {
		offDays = append(offDays, db.offDays[id])
	}
	db.mx.RUnlock()
	return offDays
}

func (db *DB) DeleteOffDay(id string) {
	db.mx.Lock()
	delete(db.offDays, id)
	db.mx.Unlock()
}
