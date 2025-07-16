package employee

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
	"unicode/utf8"
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
func (r *Repository) CreateTx(ctx context.Context, tx *sqlx.Tx, employee Entity) (int64, error) {
	var id int64
	query := "INSERT INTO employee (name) VALUES ($1) RETURNING id"
	err := tx.QueryRowContext(ctx, query, employee.Name).Scan(&id)
	return id, err
}

// найти элемент коллекции по его id
func (r *Repository) FindById(ctx context.Context, id int64) (employee Entity, err error) {
	query := "SELECT * FROM employee WHERE id=$1"
	err = r.db.GetContext(ctx, &employee, query, id)
	return employee, err
}

// найти все элементы коллекции
func (r *Repository) GetAll(ctx context.Context) (employee []Entity, err error) {
	query := "SELECT * FROM employee"
	err = r.db.SelectContext(ctx, &employee, query)
	return employee, err
}

// найти слайс элементов коллекции по слайсу их id
func (r *Repository) FindByIds(ctx context.Context, ids []int64) ([]Entity, error) {
	if len(ids) == 0 {
		return []Entity{}, fmt.Errorf("employee ids cannot be empty")
	}

	query := "SELECT * FROM employee WHERE id = ANY($1)"

	var employees []Entity
	err := r.db.SelectContext(ctx, &employees, query, pq.Int64Array(ids))

	return employees, err
}

// удалить элемент коллекции по его id
func (r *Repository) DeleteById(ctx context.Context, id int64) error {
	query := "DELETE FROM employee WHERE id=$1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// удалить элементы по слайсу их id
func (r *Repository) DeleteByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("employee ids cannot be empty")
	}

	query := "DELETE FROM employee WHERE id = ANY($1)"

	_, err := r.db.ExecContext(ctx, query, pq.Int64Array(ids))
	return err
}

// поиск сотрудника по имени
func (r *Repository) FindByName(ctx context.Context, tx *sqlx.Tx, name string) (isExists bool, err error) {
	query := "SELECT EXISTS(SELECT * FROM employee WHERE name = $1)"
	err = tx.GetContext(ctx, &isExists, query, name)
	return isExists, err
}

// FindPage возвращает сотрудников с учетом пагинации (limit, offset)
func (r *Repository) FindPage(ctx context.Context, offset int, limit int, textFilter string) ([]Entity, error) {
	var employees []Entity
	sb := strings.Builder{}
	var args []interface{}

	sb.WriteString("SELECT * FROM employee WHERE 1=1")
	if utf8.RuneCountInString(textFilter) >= 3 {
		sb.WriteString(" AND name ILIKE $1")
		args = append(args, "%"+textFilter+"%")
		sb.WriteString(" OFFSET $2 LIMIT $3")
		args = append(args, offset, limit)
	} else {
		sb.WriteString(" OFFSET $1 LIMIT $2")
		args = append(args, offset, limit)
	}

	err := r.db.SelectContext(ctx, &employees, sb.String(), args...)
	return employees, err
}

// CountAll возвращает кол-во записей
func (r *Repository) CountAll(ctx context.Context, textFilter string) (int64, error) {
	var total int64
	sb := strings.Builder{}
	var args []interface{}

	sb.WriteString("SELECT COUNT(*) FROM employee WHERE 1=1")
	textFilter = strings.TrimSpace(textFilter)
	if utf8.RuneCountInString(textFilter) >= 3 {
		sb.WriteString(" AND name ILIKE $1")
		args = append(args, "%"+textFilter+"%")
	}
	err := r.db.GetContext(ctx, &total, sb.String(), args...)
	return total, err
}
