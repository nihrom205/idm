package employee

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit"
	"github.com/jmoiron/sqlx"
	"github.com/nihrom205/idm/inner/common"
	"github.com/nihrom205/idm/inner/common/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"regexp"
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

func (m *MockRepo) CreateTx(ctx context.Context, tx *sqlx.Tx, employee Entity) (int64, error) {
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

func (m *MockRepo) FindByName(ctx context.Context, tx *sqlx.Tx, name string) (bool, error) {
	args := m.Called(tx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) BeginTransaction() (*sqlx.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (m *MockRepo) FindPage(ctx context.Context, offset int, limit int, textFilter string) ([]Entity, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) CountAll(ctx context.Context, textFilter string) (int64, error) {
	args := m.Called(textFilter)
	return args.Get(0).(int64), args.Error(1)
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

func TestCreateIfNotEmployee(t *testing.T) {
	a := assert.New(t)

	// сохранение сотрудника
	t.Run("should return id error nil", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		a.NoError(err)

		sqlxDB := sqlx.NewDb(db, "sqlmock")

		repo := Repository{db: sqlxDB}
		srv := NewService(&repo, validator.NewValidator())
		entity := getEntity()

		// Настраиваем mock для начала транзакции
		mock.ExpectBegin()

		// Настраиваем mock для проверки существования сотрудника
		mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT * FROM employee WHERE name = $1)")).
			WithArgs(entity.Name).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Настраиваем mock для создания сотрудника
		mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO employee (name) VALUES ($1) RETURNING id")).
			WithArgs(entity.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entity.Id))

		// Настраиваем mock для коммита транзакции
		mock.ExpectCommit()

		id, err := srv.Create(context.Background(), CreateRequest{Name: entity.Name})
		a.Nil(err)
		a.NotNil(id)
		a.Equal(entity.Id, id)
	})

	// не сохраняется сотрудник т.к. уже есть с таким именеи
	t.Run("should return zero error nil", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		a.NoError(err)

		sqlxDB := sqlx.NewDb(db, "sqlmock")

		repo := &Repository{db: sqlxDB}
		srv := NewService(repo, validator.NewValidator())
		entity := getEntity()

		// Настраиваем mock для начала транзакции
		mock.ExpectBegin()

		// Настраиваем mock для проверки существования сотрудника
		mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT * FROM employee WHERE name = $1)")).
			WithArgs(entity.Name).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Настраиваем mock для коммита транзакции
		mock.ExpectCommit()

		id, err := srv.Create(context.Background(), CreateRequest{Name: entity.Name})
		a.NotNil(err)
		a.NotNil(id)
		a.Equal(int64(0), id)
	})

	// работника с таким именем нет в базе данных, но создание нового работника завершилось ошибкой
	t.Run("should return error_1", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("An error occurred while creating mock: %s", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repo := &Repository{db: sqlxDB}
		srv := NewService(repo, validator.NewValidator())
		entity := getEntity()

		mock.ExpectBegin()

		// Настраиваем mock для проверки существования сотрудника
		mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT * FROM employee WHERE name = $1)")).
			WithArgs(entity.Name).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Настраиваем mock для создания сотрудника
		mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO employee (name) VALUES ($1) RETURNING id")).
			WithArgs(entity.Name).
			WillReturnError(errors.New("error insert failed"))

		id, err := srv.Create(context.Background(), CreateRequest{Name: entity.Name})
		a.Equal(int64(0), id)
		a.NotNil(err)
		a.ErrorContains(err, "error insert failed")
	})

	// не удалось создать транзакцию
	t.Run("should return error_2", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repo := &Repository{db: sqlxDB}
		service := NewService(repo, validator.NewValidator())

		mock.ExpectBegin().WillReturnError(fmt.Errorf("error create tx"))

		_, err = service.Create(context.Background(), CreateRequest{Name: getEntity().Name})
		a.NotNil(err)
		a.ErrorContains(err, "error create tx")
	})

	// ошибка при проверке наличия работника с таким именем
	t.Run("should return error_3", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repo := &Repository{db: sqlxDB}
		service := NewService(repo, validator.NewValidator())
		entity := getEntity()

		mock.ExpectBegin()

		// Настраиваем mock для проверки существования сотрудника
		mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT * FROM employee WHERE name = $1)")).
			WithArgs(entity.Name).
			WillReturnError(errors.New("error find failed"))

		id, err := service.Create(context.Background(), CreateRequest{Name: entity.Name})
		a.Equal(int64(0), id)
		a.NotNil(err)
		a.ErrorContains(err, "error find failed")
	})
}

func TestRepositoryFindPage(t *testing.T) {
	a := assert.New(t)

	t.Run("should return err validation PageSize < 1", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, validator.NewValidator())
		request := PageRequest{
			PageSize:   0,
			PageNumber: 1,
		}
		_, err := srv.FindPage(context.Background(), request)
		a.NotNil(err)
		var validateErr common.RequestValidatorError
		ok := errors.As(err, &validateErr)
		a.True(ok)
		a.Contains(validateErr.Message, "Field validation")
	})

	t.Run("should return err validation PageSize > 100", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, validator.NewValidator())
		request := PageRequest{
			PageSize:   101,
			PageNumber: 1,
		}
		_, err := srv.FindPage(context.Background(), request)
		a.NotNil(err)
		var validateErr common.RequestValidatorError
		ok := errors.As(err, &validateErr)
		a.True(ok)
		a.Contains(validateErr.Message, "Field validation")
	})

	t.Run("should return err validation PageNumber < 0", func(t *testing.T) {
		repo := &MockRepo{}
		srv := NewService(repo, validator.NewValidator())
		request := PageRequest{
			PageSize:   1,
			PageNumber: -1,
		}
		_, err := srv.FindPage(context.Background(), request)
		a.NotNil(err)
		var validateErr common.RequestValidatorError
		ok := errors.As(err, &validateErr)
		a.True(ok)
		a.Contains(validateErr.Message, "Field validation")
	})
}

func getEntity() Entity {
	return Entity{
		Id:       gofakeit.Int64(),
		Name:     gofakeit.Name(),
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
