package tests

import (
	"github.com/nihrom205/idm/inner/role"
)

type FixtureRole struct {
	role *role.Repository
}

func NewFixtureRole(role *role.Repository) *FixtureRole {
	return &FixtureRole{role}
}

func (f *FixtureRole) Role(name string) int64 {
	entity := role.Entity{
		Name: name,
	}
	newId, err := f.role.Create(entity)
	if err != nil {
		panic(err)
	}
	return newId
}
