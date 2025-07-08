package role

import (
	"context"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/nihrom205/idm/inner/common/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindById(ctx context.Context, id int64) (Entity, error) {
	args := m.Called(id)
	return args.Get(0).(Entity), args.Error(1)
}

func (m *MockRepo) Create(ctx context.Context, employee Entity) (int64, error) {
	args := m.Called(employee)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) GetAll(ctx context.Context) ([]Entity, error) {
	args := m.Called()
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) FindByIds(ctx context.Context, ids []int64) ([]Entity, error) {
	args := m.Called(ids)
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) DeleteById(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) DeleteByIds(ctx context.Context, ids []int64) error {
	args := m.Called(ids)
	return args.Error(0)
}

func TestFindById(t *testing.T) {
	a := assert.New(t)

	t.Run("should return found employee", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, nil)
		entity := getEntity()
		want := entity.toResponse()

		repo.On("FindById", int64(1)).Return(entity, nil)

		got, err := srv.FindById(context.Background(), 1)

		a.Nil(err)

		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("should return empty employee and err", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, nil)
		entity := Entity{}
		err := errors.New("database error")

		want := fmt.Errorf("error finding employee with id %d: %w", 1, err)

		repo.On("FindById", int64(1)).Return(entity, err)
		response, got := srv.FindById(context.Background(), 1)

		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})
}

func TestCreate(t *testing.T) {
	a := assert.New(t)

	t.Run("should return id", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, validator.NewValidator())
		entity := getEntity()
		request := CreateRequest{Name: entity.Name}

		repo.On("Create", request.ToEntity()).Return(int64(1), nil)
		id, err := srv.Create(context.Background(), request)

		a.Nil(err)
		a.Equal(int64(1), id)
		a.True(repo.AssertNumberOfCalls(t, "Create", 1))
	})

	t.Run("should return err", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, validator.NewValidator())
		entity := Entity{}
		request := CreateRequest{Name: entity.Name}
		err := errors.New("database error")

		repo.On("Create", request.ToEntity()).Return(int64(1), err)
		response, got := srv.Create(context.Background(), request)

		a.Empty(response)
		a.NotNil(got)
		a.True(repo.AssertNumberOfCalls(t, "Create", 0))
	})
}

func TestGetAll(t *testing.T) {
	a := assert.New(t)

	t.Run("should return all employees", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, nil)
		entities := getSliceEntity(4)

		repo.On("GetAll").Return(entities, nil)
		got, err := srv.GetAll(context.Background())

		a.Nil(err)
		a.Equal(len(entities), len(got))
		a.True(repo.AssertNumberOfCalls(t, "GetAll", 1))
	})

	t.Run("should return empty employees", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, nil)
		entities := getSliceEntity(0)
		err := errors.New("database error")

		want := fmt.Errorf("error getting all employees: %w", err)

		repo.On("GetAll").Return(entities, err)
		response, got := srv.GetAll(context.Background())

		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "GetAll", 1))
	})
}

func TestFindByIds(t *testing.T) {
	a := assert.New(t)

	t.Run("should return employees", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, nil)
		entities := getSliceEntity(3)

		findByIds := []int64{entities[0].Id, entities[1].Id, entities[2].Id}

		repo.On("FindByIds", mock.Anything).Return(entities, nil)
		got, err := srv.FindByIds(context.Background(), findByIds)

		a.Nil(err)
		a.Equal(len(entities), len(got))
		a.True(repo.AssertNumberOfCalls(t, "FindByIds", 1))
	})

	t.Run("should return empty employee", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, nil)
		entities := getSliceEntity(0)

		err := errors.New("database error")

		findByIds := []int64{1, 2, 3}
		want := fmt.Errorf("error finding employee with id %d: %w", findByIds, err)

		repo.On("FindByIds", mock.Anything).Return(entities, err)
		response, got := srv.FindByIds(context.Background(), findByIds)

		a.Empty(response)
		a.NotNil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindByIds", 1))
	})
}

func TestDeleteById(t *testing.T) {
	a := assert.New(t)

	t.Run("should return error nil", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, nil)
		deleteById := int64(1)

		repo.On("DeleteById", deleteById).Return(nil)
		err := srv.DeleteById(context.Background(), deleteById)

		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})

	t.Run("should return error", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, nil)
		deleteById := int64(1)

		err := errors.New("database error")

		want := fmt.Errorf("error deleting employee with id %d: %w", deleteById, err)

		repo.On("DeleteById", deleteById).Return(err)
		got := srv.DeleteById(context.Background(), deleteById)

		a.NotNil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})
}

func TestDeleteByIds(t *testing.T) {
	a := assert.New(t)

	t.Run("should return error nil", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, nil)
		deleteByIds := []int64{1, 2, 3}

		repo.On("DeleteByIds", deleteByIds).Return(nil)
		err := srv.DeleteByIds(context.Background(), deleteByIds)

		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteByIds", 1))
	})

	t.Run("should return error nil", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, nil)
		deleteByIds := []int64{1, 2, 3}

		err := errors.New("database error")

		want := fmt.Errorf("error deleting employee with id %d: %w", deleteByIds, err)

		repo.On("DeleteByIds", deleteByIds).Return(err)
		got := srv.DeleteByIds(context.Background(), deleteByIds)

		a.NotNil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteByIds", 1))
	})
}

func getEntity() Entity {
	roles := []string{"admin", "user", "manager", "guest"}
	return Entity{
		Id:       gofakeit.Int64(),
		Name:     gofakeit.RandString(roles),
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
}

func getSliceEntity(countItems int) []Entity {
	entities := make([]Entity, 0, countItems)
	for i := 0; i < countItems; i++ {
		entities = append(entities, getEntity())
	}
	return entities
}
