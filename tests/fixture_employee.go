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

	tx, err := f.employee.BeginTransaction()
	if err != nil {
		panic(err)
	}
	newId, err := f.employee.CreateTx(tx, entity)
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		return 0
	}
	if err != nil {
		panic(err)
	}
	return newId
}
