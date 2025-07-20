package role

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nihrom205/idm/inner/cache"
	"time"
)

const ROLE = "role"

type Repository struct {
	cache *cache.RedisCache
	db    *sqlx.DB
}

func NewRoleRepository(cache *cache.RedisCache, db *sqlx.DB) *Repository {
	return &Repository{
		cache: cache,
		db:    db,
	}
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
	cacheKey := r.cache.GetCacheKey(ROLE, id)

	// Пробуем получить из кэша
	if found := r.cache.Get(ctx, cacheKey, &role); found {
		return role, nil
	}

	// получаем из БД
	query := "SELECT * FROM role WHERE id=$1"
	err = r.db.GetContext(ctx, &role, query, id)
	if err == nil {
		// сохраняем в cache
		errCache := r.cache.Set(ctx, cacheKey, role, time.Minute*5)
		if errCache != nil {
			log.Errorf("error caching role: %v", errCache.Error())
		}
	}

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

	errCache := r.cache.InvalidateCache(ctx, r.cache.GetCacheKey(ROLE, id))
	if errCache != nil {
		log.Errorf("error caching employee: %v", errCache.Error())
	}

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
