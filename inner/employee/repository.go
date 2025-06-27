package employee

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// запрос транзакции у БД
func (r *Repository) BeginTransaction() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

// добавить новый элемент в коллекцию
func (r *Repository) CreateTx(tx *sqlx.Tx, employee Entity) (int64, error) {
	var id int64
	query := "INSERT INTO employee (name) VALUES ($1) RETURNING id"
	err := tx.QueryRow(query, employee.Name).Scan(&id)
	return id, err
}

// найти элемент коллекции по его id
func (r *Repository) FindById(id int64) (employee Entity, err error) {
	query := "SELECT * FROM employee WHERE id=$1"
	err = r.db.Get(&employee, query, id)
	return employee, err
}

// найти все элементы коллекции
func (r *Repository) GetAll() (employee []Entity, err error) {
	query := "SELECT * FROM employee"
	err = r.db.Select(&employee, query)
	return employee, err
}

// найти слайс элементов коллекции по слайсу их id
func (r *Repository) FindByIds(ids []int64) ([]Entity, error) {
	if len(ids) == 0 {
		return []Entity{}, fmt.Errorf("employee ids cannot be empty")
	}

	query := "SELECT * FROM employee WHERE id = ANY($1)"

	var employees []Entity
	err := r.db.Select(&employees, query, pq.Int64Array(ids))

	return employees, err
}

// удалить элемент коллекции по его id
func (r *Repository) DeleteById(id int64) error {
	query := "DELETE FROM employee WHERE id=$1"
	_, err := r.db.Exec(query, id)
	return err
}

// удалить элементы по слайсу их id
func (r *Repository) DeleteByIds(ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("employee ids cannot be empty")
	}

	query := "DELETE FROM employee WHERE id = ANY($1)"

	_, err := r.db.Exec(query, pq.Int64Array(ids))
	return err
}

// поиск сотрудника по имени
func (r *Repository) FindByName(tx *sqlx.Tx, name string) (isExists bool, err error) {
	query := "SELECT EXISTS(SELECT * FROM employee WHERE name = $1)"
	err = tx.Get(&isExists, query, name)
	return isExists, err
}
