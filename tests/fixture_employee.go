package tests

import (
	"github.com/nihrom205/idm/inner/employee"
)

type FixtureEmployee struct {
	employee *employee.Repository
}

func NewFixtureEmployee(employee *employee.Repository) *FixtureEmployee {
	return &FixtureEmployee{employee}
}

func (f *FixtureEmployee) Employee(name string) int64 {
	entity := employee.Entity{
		Name: name,
	}

	newId, err := f.employee.Create(entity)
	if err != nil {
		panic(err)
	}
	return newId
}
