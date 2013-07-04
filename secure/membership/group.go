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

type Group struct {
	Id        model.Identifier `bson:"-" datastore:"-"`
	Name      string
	Info      GroupInfo
	Privilege map[string]bool
}

func (g *Group) GetId() model.Identifier {
	return g.Id
}

func (g *Group) SetId(id model.Identifier) error {
	g.Id = id
	return nil
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
