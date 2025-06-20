package employee

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type StubRepo struct{}

func (s *StubRepo) FindById(id int64) (Entity, error) {
	return Entity{
		Id:       id,
		Name:     "Test User",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}, nil
}

func (s *StubRepo) Create(employee Entity) (int64, error) {
	return 0, nil
}

func (s *StubRepo) GetAll() ([]Entity, error) {
	return []Entity{}, nil
}

func (s *StubRepo) FindByIds(ids []int64) ([]Entity, error) {
	return []Entity{}, nil
}

func (s *StubRepo) DeleteById(id int64) error {
	return nil
}

func (s *StubRepo) DeleteByIds(ids []int64) error {
	return nil
}

func TestStubFindById(t *testing.T) {
	a := assert.New(t)

	t.Run("should return found employee", func(t *testing.T) {
		repo := &StubRepo{}
		srv := NewService(repo)

		got, err := srv.FindById(1)

		a.Nil(err)

		a.Equal(int64(1), got.Id)
		a.Equal("Test User", got.Name)
	})

}
