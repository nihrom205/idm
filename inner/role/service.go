package role

import "fmt"

type Repo interface {
	Create(role Entity) (int64, error)
	FindById(id int64) (Entity, error)
	GetAll() (role []Entity, err error)
	FindByIds(ids []int64) ([]Entity, error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(role Entity) (int64, error) {
	id, err := s.repo.Create(role)
	if err != nil {
		return 0, fmt.Errorf("error failed to create employee with id %d: %w", id, err)
	}
	return id, nil
}

func (s *Service) FindById(id int64) (Response, error) {
	role, err := s.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee with id %d: %w", id, err)
	}

	return role.toResponse(), nil
}

func (s *Service) GetAll() ([]Response, error) {
	roles, err := s.repo.GetAll()
	if err != nil {
		return []Response{}, fmt.Errorf("error getting all employees: %w", err)
	}

	response := make([]Response, 0, len(roles))
	for _, item := range roles {
		response = append(response, item.toResponse())
	}

	return response, nil
}

func (s *Service) FindByIds(ids []int64) ([]Response, error) {
	roles, err := s.repo.FindByIds(ids)
	if err != nil {
		return []Response{}, fmt.Errorf("error finding employee with id %d: %w", ids, err)
	}

	response := make([]Response, 0, len(roles))
	for _, item := range roles {
		response = append(response, item.toResponse())
	}

	return response, nil
}

func (s *Service) DeleteById(id int64) error {
	err := s.repo.DeleteById(id)
	if err != nil {
		return fmt.Errorf("error deleting employee with id %d: %w", id, err)
	}

	return nil
}

func (s *Service) DeleteByIds(ids []int64) error {
	err := s.repo.DeleteByIds(ids)
	if err != nil {
		return fmt.Errorf("error deleting employee with id %d: %w", ids, err)
	}

	return nil
}
