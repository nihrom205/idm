package employee

type CreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=155"`
}

func (r *CreateRequest) ToEntity() Entity {
	return Entity{Name: r.Name}
}

type FindByIdRequest struct {
	Id int64 `json:"id" validate:"required,gt=0"`
}

type FindEmployeesByIdsRequest struct {
	Ids []int64 `json:"ids" validate:"required,min=1,dive,gt=0"`
}

type DeleteEmployeesByIdsRequest struct {
	Ids []int64 `json:"ids" validate:"required,min=1,dive,gt=0"`
}
