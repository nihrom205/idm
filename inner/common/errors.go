package common

type RequestValidatorError struct {
	Message string
}

func (e RequestValidatorError) Error() string {
	return e.Message
}

type AlreadyExistsError struct {
	Message string
}

func (e AlreadyExistsError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

type RepositoryError struct {
	Message string
}

func (e RepositoryError) Error() string {
	return e.Message
}
