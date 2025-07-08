package role

import (
	"context"
	"fmt"
	"github.com/nihrom205/idm/inner/common"
)

type Repo interface {
	Create(ctx context.Context, role Entity) (int64, error)
	FindById(ctx context.Context, id int64) (Entity, error)
	GetAll(ctx context.Context) (role []Entity, err error)
	FindByIds(ctx context.Context, ids []int64) ([]Entity, error)
	DeleteById(ctx context.Context, id int64) error
	DeleteByIds(ctx context.Context, ids []int64) error
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

func (s *Service) Create(ctx context.Context, request CreateRequest) (int64, error) {

	// валидируем запрос
	err := s.validator.Validate(request)
	if err != nil {
		// возвращаем кастомную ошибку в случае, если запрос не прошёл валидацию
		return 0, common.RequestValidatorError{Message: err.Error()}
	}
	id, err := s.repo.Create(ctx, request.ToEntity())
	if err != nil {
		return 0, fmt.Errorf("error failed to create employee with id %d: %w", id, err)
	}
	return id, nil
}

func (s *Service) FindById(ctx context.Context, id int64) (Response, error) {
	role, err := s.repo.FindById(ctx, id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee with id %d: %w", id, err)
	}

	return role.toResponse(), nil
}

func (s *Service) GetAll(ctx context.Context) ([]Response, error) {
	roles, err := s.repo.GetAll(ctx)
	if err != nil {
		return []Response{}, fmt.Errorf("error getting all employees: %w", err)
	}

	response := make([]Response, 0, len(roles))
	for _, item := range roles {
		response = append(response, item.toResponse())
	}

	return response, nil
}

func (s *Service) FindByIds(ctx context.Context, ids []int64) ([]Response, error) {
	roles, err := s.repo.FindByIds(ctx, ids)
	if err != nil {
		return []Response{}, fmt.Errorf("error finding employee with id %d: %w", ids, err)
	}

	response := make([]Response, 0, len(roles))
	for _, item := range roles {
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
