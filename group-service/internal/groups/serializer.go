package groups

type Member struct {
	ID   string   `json:"id" validate:"required"`
	Role string `json:"role" validate:"required"`
}

type GroupCreateSerizliser struct {
	Name    string   `json:"name" validate:"required,min=3,max=20"`
	Members []Member `json:"members" validate:"required,dive"`
}

type GroupGetSerizliser struct {
	ID   string   `param:"id" validate:"required"`
}
