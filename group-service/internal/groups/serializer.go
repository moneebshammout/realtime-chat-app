package groups

type Member struct {
	ID   uint   `json:"id" validate:"required"`
	Role string `json:"role" validate:"required"`
}

type GroupSerizliser struct {
	Name    string   `json:"name" validate:"required,min=3,max=20"`
	Members []Member `json:"members" validate:"required,dive"`
}
