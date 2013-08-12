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
	At time.Time
}

func (s *SessionInfo) GetCreateTime() time.Time {
	return s.At
}

func (s *SessionInfo) SetCreateTime(t time.Time) {
	s.At = t
}
