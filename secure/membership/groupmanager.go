package membership

import (
	"github.com/kidstuff/toys/model"
)

type GroupManager interface {
	AddGroupDetail(name string, info *GroupInfo, pri map[string]bool) error
	UpdateInfo(id model.Identifier, info *GroupInfo) error
	UpdatePrivilege(id model.Identifier, pri map[string]bool) error
	FindGroup(id model.Identifier) (Grouper, error)
	FindAllGroup(offsetId model.Identifier, limit int) ([]Grouper, error)
}
