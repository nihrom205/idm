package employee

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type StubRepo struct{}

func (s *StubRepo) FindById(ctx context.Context, id int64) (Entity, error) {
	return Entity{
		Id:       id,
		Name:     "Test User",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}, nil
}

func (s *StubRepo) CreateTx(ctx context.Context, tx *sqlx.Tx, employee Entity) (int64, error) {
	return 0, nil
}

func (s *StubRepo) GetAll(ctx context.Context) ([]Entity, error) {
	return []Entity{}, nil
}

func (s *StubRepo) FindByIds(ctx context.Context, ids []int64) ([]Entity, error) {
	return []Entity{}, nil
}

func (s *StubRepo) DeleteById(ctx context.Context, id int64) error {
	return nil
}

func (s *StubRepo) DeleteByIds(ctx context.Context, ids []int64) error {
	return nil
}

func (s *StubRepo) FindByName(ctx context.Context, tx *sqlx.Tx, name string) (bool, error) {
	return false, nil
}

func (s *StubRepo) BeginTransaction() (*sqlx.Tx, error) {
	return nil, nil
}

func (s *StubRepo) FindPage(ctx context.Context, offset int, limit int, textFilter string) ([]Entity, error) {
	return []Entity{}, nil
}

func (s *StubRepo) CountAll(ctx context.Context, textFilter string) (int64, error) {
	return 0, nil
}

func TestStubFindById(t *testing.T) {
	a := assert.New(t)

	t.Run("should return found employee", func(t *testing.T) {
		repo := &StubRepo{}
		srv := NewService(repo, nil)

		got, err := srv.FindById(context.Background(), 1)

		a.Nil(err)

		a.Equal(int64(1), got.Id)
		a.Equal("Test User", got.Name)
	})

}
