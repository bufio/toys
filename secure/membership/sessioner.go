package membership

import (
	"github.com/kidstuff/toys/model"
	"time"
)

type Sessioner interface {
	GetId() model.Identifier
	SetId(model.Identifier) error
	GetCreateTime() time.Time
	SetCreateTime(time.Time)
}

type SessionInfo struct {
	Id model.Identifier `bson:"-" datastore:"-"`
	At time.Time
}

var _ Sessioner = &SessionInfo{}

// GetId just an virtual function, you may want to re-implement it
func (s *SessionInfo) GetId() model.Identifier {
	return s.Id
}

// SetId just an virtual function, you may want to re-implement it
func (s *SessionInfo) SetId(id model.Identifier) error {
	s.Id = id
	return nil
}

func (s *SessionInfo) GetCreateTime() time.Time {
	return s.At
}

func (s *SessionInfo) SetCreateTime(t time.Time) {
	s.At = t
}
