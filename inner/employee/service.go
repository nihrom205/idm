package employee

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nihrom205/idm/inner/common"
)

type Repo interface {
	CreateTx(ctx context.Context, tx *sqlx.Tx, employee Entity) (int64, error)
	FindById(ctx context.Context, id int64) (Entity, error)
	GetAll(ctx context.Context) (employee []Entity, err error)
	FindByIds(ctx context.Context, ids []int64) ([]Entity, error)
	DeleteById(ctx context.Context, id int64) error
	DeleteByIds(ctx context.Context, ids []int64) error
	FindByName(ctx context.Context, tx *sqlx.Tx, name string) (bool, error)
	BeginTransaction() (*sqlx.Tx, error)
}

type Validator interface {
	Validate(request any) error
}

type Service struct {
	repo      Repo
	validator Validator
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{
		repo:      repo,
		validator: validator,
	}
}

// Метод для создания нового сотрудника
// принимает на вход CreateRequest - структура запроса на создание сотрудника
func (s *Service) Create(ctx context.Context, request CreateRequest) (int64, error) {

	// валидируем запрос
	err := s.validator.Validate(request)
	if err != nil {
		// возвращаем кастомную ошибку в случае, если запрос не прошёл валидацию
		return 0, common.RequestValidatorError{Message: err.Error()}
	}

	tx, err := s.repo.BeginTransaction()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("creating employee panic: %v", r)
			// если была паника, то откатываем транзакцию
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else if err != nil {
			// если произошла другая ошибка (не паника), то откатываем транзакцию
			if tx != nil {
				errTx := tx.Rollback()
				if errTx != nil {
					err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
				}
			}
		} else {
			// если ошибок нет, то коммитим транзакцию
			errTx := tx.Commit()
			if errTx != nil {
				err = fmt.Errorf("creating employee: commiting transaction error: %w", errTx)
			}
		}
	}()

	if err != nil {
		return int64(0), fmt.Errorf("error creating transaction: %w", err)
	}

	// в рамках транзакции проверяем наличие в базе данных работника с таким же именем
	isExist, err := s.repo.FindByName(ctx, tx, request.Name)
	if err != nil {
		return 0, fmt.Errorf("error finding employee by name: %s, %w", request.Name, err)
	}
	if isExist {
		return 0, common.AlreadyExistsError{Message: fmt.Sprintf("employee with name %s already exists", request.Name)}
	}

	// в случае отсутствия сотрудника с таким же именем - в рамках этой же транзакции вызываем метод репозитория,
	// который должен будет создать нового сотрудника
	newEmployeeId, err := s.repo.CreateTx(ctx, tx, request.ToEntity())
	if err != nil {
		return 0, fmt.Errorf("error failed to create employee with id %d: %w", newEmployeeId, err)
	}

	return newEmployeeId, nil
}

func (s *Service) FindById(ctx context.Context, id int64) (Response, error) {
	employees, err := s.repo.FindById(ctx, id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee with id %d: %w", id, err)
	}

	return employees.toResponse(), nil
}

func (s *Service) GetAll(ctx context.Context) ([]Response, error) {
	employees, err := s.repo.GetAll(ctx)
	if err != nil {
		return []Response{}, fmt.Errorf("error getting all employees: %w", err)
	}

	response := make([]Response, 0, len(employees))
	for _, item := range employees {
		response = append(response, item.toResponse())
	}

	return response, nil
}

func (s *Service) FindByIds(ctx context.Context, ids []int64) ([]Response, error) {
	employee, err := s.repo.FindByIds(ctx, ids)
	if err != nil {
		return []Response{}, fmt.Errorf("error finding employee with id %d: %w", ids, err)
	}

	response := make([]Response, 0, len(employee))
	for _, item := range employee {
		response = append(response, item.toResponse())
	}

	return response, nil
}

func (s *Service) DeleteById(ctx context.Context, id int64) error {
	err := s.repo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting employee with id %d: %w", id, err)
	}

	return nil
}

func (s *Service) DeleteByIds(ctx context.Context, ids []int64) error {
	err := s.repo.DeleteByIds(ctx, ids)
	if err != nil {
		return fmt.Errorf("error deleting employee with id %d: %w", ids, err)
	}

	return nil
}
