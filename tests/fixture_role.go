package tests

import (
	"github.com/nihrom205/idm/inner/role"
)

type FixtureRole struct {
	role *role.RoleRepository
}

func NewFixtureRole(role *role.RoleRepository) *FixtureRole {
	return &FixtureRole{role}
}

func (f *FixtureRole) Role(name string) int64 {
	entity := role.RoleEntity{
		Name: name,
	}
	newId, err := f.role.Create(entity)
	if err != nil {
		panic(err)
	}
	return newId
}
