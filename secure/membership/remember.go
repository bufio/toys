package membership

import (
	"github.com/kidstuff/toys/model"
	"time"
)

type Remember interface {
	GetId() model.Identifier
	SetId(interface{}) error
	GetToken() string
	SetToken(string)
	GetExpiration() time.Time
	SetExpiration(time.Time)
}

type RememberInfo struct {
	Id    model.Identifier `bson:"-" datastore:"-"`
	Token string
	Exp   time.Time
}

// GetId just an virtual function, you may want to re-implement it
func (r *RememberInfo) GetId() model.Identifier {
	return r.Id
}

// SetId just an virtual function, you may want to re-implement it
func (r *RememberInfo) SetId(id model.Identifier) error {
	r.Id = id
	return nil
}

func (r *RememberInfo) GetToken() string {
	return r.Token
}

func (r *RememberInfo) SetToken(token string) {
	r.Token = token
}

func (r *RememberInfo) GetExpiration() time.Time {
	return r.Exp
}

func (r *RememberInfo) SetExpiration(t time.Time) {
	r.Exp = t
}
