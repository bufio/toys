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
	Token string
	Exp   time.Time
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
