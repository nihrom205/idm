package role

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// добавить новый элемент в коллекцию
func (r *Repository) Create(ctx context.Context, role Entity) (int64, error) {
	var id int64
	query := "INSERT INTO role (name) VALUES ($1) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, role.Name).Scan(&id)
	return id, err
}

// найти элемент коллекции по его id
func (r *Repository) FindById(ctx context.Context, id int64) (role Entity, err error) {
	query := "SELECT * FROM role WHERE id=$1"
	err = r.db.GetContext(ctx, &role, query, id)
	return role, err
}

// найти все элементы коллекции
func (r *Repository) GetAll(ctx context.Context) (roles []Entity, err error) {
	query := "SELECT * FROM role"
	err = r.db.SelectContext(ctx, &roles, query)
	return roles, err
}

// найти слайс элементов коллекции по слайсу их id
func (r *Repository) FindByIds(ctx context.Context, ids []int64) ([]Entity, error) {
	if len(ids) == 0 {
		return []Entity{}, fmt.Errorf("role ids cannot be empty")
	}

	query := "SELECT * FROM role WHERE id = ANY($1)"

	var roles []Entity
	err := r.db.SelectContext(ctx, &roles, query, pq.Int64Array(ids))
	return roles, err
}

// удалить элемент коллекции по его id
func (r *Repository) DeleteById(ctx context.Context, id int64) error {
	query := "DELETE FROM role WHERE id=$1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// удалить элементы по слайсу их id
func (r *Repository) DeleteByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("role ids cannot be empty")
	}

	query := "DELETE FROM role WHERE id = ANY($1)"

	_, err := r.db.ExecContext(ctx, query, pq.Int64Array(ids))
	return err
}
