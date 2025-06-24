package tests

import (
	"github.com/jmoiron/sqlx"
	"github.com/nihrom205/idm/inner/database"
	"github.com/nihrom205/idm/inner/employee"
	"github.com/nihrom205/idm/inner/role"
	"github.com/stretchr/testify/assert"
	"testing"
)

func clearDatabaseEmployee(db *sqlx.DB) {
	db.MustExec("DELETE FROM employee")
}

func createTableEmployee(db *sqlx.DB) {
	query := `CREATE TABLE IF NOT EXISTS  employee (
    id bigint generated always as IDENTITY primary key not null,
    name text not null,
    create_at timestamptz default now(),
    update_at timestamptz default now())`

	db.MustExec(query)
}

func clearDatabaseRole(db *sqlx.DB) {
	db.MustExec("DELETE FROM role")
}

func createTableRole(db *sqlx.DB) {
	query := `CREATE TABLE IF NOT EXISTS  role (
    id bigint generated always as IDENTITY primary key not null,
    name text not null unique,
    create_at timestamptz default now(),
    update_at timestamptz default now())`

	db.MustExec(query)
}

func TestRepositoryEmployee(t *testing.T) {
	a := assert.New(t)
	db := database.Connect()
	createTableEmployee(db)

	defer func() {
		if r := recover(); r != nil {
			clearDatabaseEmployee(db)
		}
	}()

	employeeRepository := employee.NewEmployeeRepository(db)

	fixture := NewFixtureEmployee(employeeRepository)

	t.Run("FindByIds", func(t *testing.T) {
		newEmployeeId := fixture.Employee("TestName")

		got, err := employeeRepository.FindById(newEmployeeId)
		a.Nil(err)
		a.NotEmpty(got)
		a.NotEmpty(got.Id)
		a.NotEmpty(got.CreateAt)
		a.NotEmpty(got.UpdateAt)
		a.Equal("TestName", got.Name)
		clearDatabaseEmployee(db)
	})

	t.Run("GetAll", func(t *testing.T) {
		employeeIvanId := fixture.Employee("Ivan")
		employeeVadimId := fixture.Employee("Vadim")

		got, err := employeeRepository.GetAll()
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(2, len(got))
		a.NotEmpty(got[0].Id)
		a.Equal(got[0].Id, employeeIvanId)
		a.Equal(got[0].Name, "Ivan")
		a.NotEmpty(got[1].Id)
		a.Equal(got[1].Id, employeeVadimId)
		a.Equal(got[1].Name, "Vadim")
		clearDatabaseEmployee(db)
	})

	t.Run("GetEmployeeByIds", func(t *testing.T) {
		ivanId := fixture.Employee("Ivan")
		vadimId := fixture.Employee("Vadim")
		_ = fixture.Employee("Pavel")

		ids := []int64{ivanId, vadimId}
		got, err := employeeRepository.FindByIds(ids)
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(2, len(got))
		a.NotEmpty(got[0].Id)
		a.Equal(got[0].Id, ivanId)
		a.Equal(got[0].Name, "Ivan")
		a.NotEmpty(got[1].Id)
		a.Equal(got[1].Id, vadimId)
		a.Equal(got[1].Name, "Vadim")
		clearDatabaseEmployee(db)
	})

	t.Run("DeleteEmployeeById", func(t *testing.T) {
		ivanId := fixture.Employee("Ivan")

		err := employeeRepository.DeleteById(ivanId)
		a.Nil(err)

		_, err = employeeRepository.FindById(ivanId)
		a.NotNil(err)
		a.Equal(err.Error(), "sql: no rows in result set")

		clearDatabaseEmployee(db)
	})

	t.Run("DeleteEmployeeByIds", func(t *testing.T) {
		ivanId := fixture.Employee("Ivan")
		vadimId := fixture.Employee("Vadim")
		pavelId := fixture.Employee("Pavel")

		delIds := []int64{ivanId, vadimId}
		err := employeeRepository.DeleteByIds(delIds)
		a.Nil(err)

		findIds := []int64{pavelId, vadimId, pavelId}
		got, err := employeeRepository.FindByIds(findIds)
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(1, len(got))
		a.NotEmpty(got[0].Id)
		a.Equal(got[0].Id, pavelId)
		a.Equal(got[0].Name, "Pavel")
		clearDatabaseEmployee(db)
	})

	t.Run("BeginTransaction", func(t *testing.T) {
		//_ = fixture.Employee("Ivan")
		employeeName := "Ivan"

		//delIds := []int64{ivanId, vadimId}
		tx, err := employeeRepository.BeginTransaction()
		a.Nil(err)
		var isExists bool
		query := "SELECT EXISTS(SELECT * FROM employee WHERE name = $1)"
		err = tx.Get(&isExists, query, employeeName)
		a.Nil(err)
		a.False(isExists)

		query = "INSERT INTO employee (name) VALUES ($1) RETURNING id"
		var id int64
		err = tx.QueryRow(query, employeeName).Scan(&id)
		a.Nil(err)
		a.NotEqual(int64(0), id)
		err = tx.Commit()
		if err != nil {
			return
		}
		clearDatabaseEmployee(db)
	})
}

func TestRepositoryRole(t *testing.T) {
	a := assert.New(t)
	db := database.Connect()
	createTableRole(db)

	defer func() {
		if r := recover(); r != nil {
			clearDatabaseRole(db)
		}
	}()

	roleRepository := role.NewRoleRepository(db)

	fixture := NewFixtureRole(roleRepository)

	t.Run("FindById", func(t *testing.T) {
		newRoleId := fixture.Role("Director")

		got, err := roleRepository.FindById(newRoleId)
		a.Nil(err)
		a.NotEmpty(got)
		a.NotEmpty(got.Id)
		a.NotEmpty(got.CreateAt)
		a.NotEmpty(got.UpdateAt)
		a.Equal("Director", got.Name)
		clearDatabaseRole(db)
	})

	t.Run("GetAll", func(t *testing.T) {
		managerId := fixture.Role("Manager")
		dirId := fixture.Role("Director")

		got, err := roleRepository.GetAll()
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(2, len(got))
		a.NotEmpty(got[0].Id)
		a.Equal(got[0].Id, managerId)
		a.Equal(got[0].Name, "Manager")
		a.NotEmpty(got[1].Id)
		a.Equal(got[1].Id, dirId)
		a.Equal(got[1].Name, "Director")
		clearDatabaseRole(db)
	})

	t.Run("GetEmployeeByIds", func(t *testing.T) {
		managerId := fixture.Role("Manager")
		dirId := fixture.Role("Director")
		_ = fixture.Role("Driver")

		ids := []int64{managerId, dirId}
		got, err := roleRepository.FindByIds(ids)
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(2, len(got))
		a.NotEmpty(got[0].Id)
		a.Equal(got[0].Id, managerId)
		a.Equal(got[0].Name, "Manager")
		a.NotEmpty(got[1].Id)
		a.Equal(got[1].Id, dirId)
		a.Equal(got[1].Name, "Director")
		clearDatabaseRole(db)
	})

	t.Run("DeleteEmployeeById", func(t *testing.T) {
		managerId := fixture.Role("Manager")

		err := roleRepository.DeleteById(managerId)
		a.Nil(err)

		_, err = roleRepository.FindById(managerId)
		a.NotNil(err)
		a.Equal(err.Error(), "sql: no rows in result set")

		clearDatabaseRole(db)
	})

	t.Run("DeleteEmployeeByIds", func(t *testing.T) {
		managerId := fixture.Role("Manager")
		dirId := fixture.Role("Director")
		driverId := fixture.Role("Driver")

		delIds := []int64{managerId, dirId}
		err := roleRepository.DeleteByIds(delIds)
		a.Nil(err)

		findIds := []int64{managerId, dirId, driverId}
		got, err := roleRepository.FindByIds(findIds)
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(1, len(got))
		a.NotEmpty(got[0].Id)
		a.Equal(got[0].Id, driverId)
		a.Equal(got[0].Name, "Driver")
		clearDatabaseRole(db)
	})
}
