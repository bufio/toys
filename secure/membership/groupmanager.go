package membership

import (
	"errors"
	"github.com/kidstuff/toys/model"
)

var (
	ErrDuplicateName = errors.New("membership: duplicate Group Name")
)

type GroupManager interface {
	AddGroupDetail(name string, info GroupInfo, pri map[string]bool) (Grouper, error)
	UpdateInfo(id model.Identifier, info GroupInfo) error
	UpdatePrivilege(id model.Identifier, pri map[string]bool) error
	FindGroup(id model.Identifier) (Grouper, error)
	FindGroupByName(name string) (Grouper, error)
	FindAllGroup(offsetId model.Identifier, limit int) ([]Grouper, error)
}
