package employee

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type EmployeeEntity struct {
	Id       int64     `db:"id"`
	Name     string    `db:"name"`
	CreateAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

type EmployeeRepository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(db *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

// добавить новый элемент в коллекцию
func (r *EmployeeRepository) Create(employee EmployeeEntity) (int64, error) {
	var id int64
	query := "INSERT INTO employee (name) VALUES ($1) RETURNING id"
	err := r.db.QueryRow(query, employee.Name).
		Scan(&id)
	return id, err
}

// найти элемент коллекции по его id
func (r *EmployeeRepository) FindById(id int64) (employee EmployeeEntity, err error) {
	query := "SELECT * FROM employee WHERE id=$1"
	err = r.db.Get(&employee, query, id)
	return employee, err
}

// найти все элементы коллекции
func (r *EmployeeRepository) GetAll() (employee []EmployeeEntity, err error) {
	query := "SELECT * FROM employee"
	err = r.db.Select(&employee, query)
	return employee, err
}

// найти слайс элементов коллекции по слайсу их id
func (r *EmployeeRepository) FindByIds(ids []int64) []EmployeeEntity {
	if len(ids) == 0 {
		return []EmployeeEntity{}
	}

	query := "SELECT * FROM employee WHERE id = ANY($1)"

	var employees []EmployeeEntity
	err := r.db.Select(&employees, query, pq.Int64Array(ids))
	if err != nil {
		return []EmployeeEntity{}
	}
	return employees
}

// удалить элемент коллекции по его id
func (r *EmployeeRepository) DeleteById(id int64) error {
	query := "DELETE FROM employee WHERE id=$1"
	_, err := r.db.Exec(query, id)
	return err
}

// удалить элементы по слайсу их id
func (r *EmployeeRepository) DeleteByIds(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	query := "DELETE FROM employee WHERE id = ANY($1)"

	_, err := r.db.Exec(query, pq.Int64Array(ids))
	if err != nil {
		return err
	}
	return nil
}
