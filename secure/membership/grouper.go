package membership

import (
	"github.com/kidstuff/toys/model"
)

type Grouper interface {
	GetId() model.Identifier
	SetId(model.Identifier) error
	GetName() string
	GetInfomation() GroupInfo
	GetPrivilege() map[string]bool
}

type BriefGroup struct {
	Id   model.Identifier `bson:"-" datastore:"-"`
	Name string
}

type Group struct {
	Name      string
	Info      GroupInfo `datastore:",noindex"`
	Privilege map[string]bool
}

func (g *Group) GetName() string {
	return g.Name
}

func (g *Group) GetInfomation() GroupInfo {
	return g.Info
}

func (g *Group) GetPrivilege() map[string]bool {
	return g.Privilege
}

type GroupInfo struct {
	Description string
}
