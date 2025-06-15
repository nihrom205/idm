package tests

import (
	"github.com/nihrom205/idm/inner/employee"
)

type FixtureEmployee struct {
	employee *employee.EmployeeRepository
}

func NewFixtureEmployee(employee *employee.EmployeeRepository) *FixtureEmployee {
	return &FixtureEmployee{employee}
}

func (f *FixtureEmployee) Employee(name string) int64 {
	entity := employee.EmployeeEntity{
		Name: name,
	}

	newId, err := f.employee.Create(entity)
	if err != nil {
		panic(err)
	}
	return newId
}
