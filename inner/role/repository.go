package role

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type RoleEntity struct {
	Id       int64     `db:"id"`
	Name     string    `db:"name"`
	CreateAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// добавить новый элемент в коллекцию
func (r *RoleRepository) Create(role RoleEntity) (int64, error) {
	var id int64
	query := "INSERT INTO role (name) VALUES ($1) RETURNING id"
	err := r.db.QueryRow(query, role.Name).Scan(&id)
	return id, err
}

// найти элемент коллекции по его id
func (r *RoleRepository) FindById(id int64) (role RoleEntity, err error) {
	query := "SELECT * FROM role WHERE id=$1"
	err = r.db.Get(&role, query, id)
	return role, err
}

// найти все элементы коллекции
func (r *RoleRepository) GetAll() (roles []RoleEntity, err error) {
	query := "SELECT * FROM role"
	err = r.db.Select(&roles, query)
	return roles, err
}

// найти слайс элементов коллекции по слайсу их id
func (r *RoleRepository) FindByIds(ids []int64) []RoleEntity {
	if len(ids) == 0 {
		return []RoleEntity{}
	}

	query := "SELECT * FROM role WHERE id = ANY($1)"

	var roles []RoleEntity
	err := r.db.Select(&roles, query, pq.Int64Array(ids))
	if err != nil {
		return []RoleEntity{}
	}
	return roles
}

// удалить элемент коллекции по его id
func (r *RoleRepository) DeleteById(id int64) error {
	query := "DELETE FROM role WHERE id=$1"
	_, err := r.db.Exec(query, id)
	return err
}

// удалить элементы по слайсу их id
func (r *RoleRepository) DeleteByIds(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	query := "DELETE FROM role WHERE id = ANY($1)"

	_, err := r.db.Exec(query, pq.Int64Array(ids))
	if err != nil {
		return err
	}
	return nil
}
