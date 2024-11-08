package groups

type IDParam struct {
	ID uint `param:"id" validate:"required"`
}

type Member struct {
	ID   string `json:"id" validate:"required"`
	Role string `json:"role" validate:"required"`
}

type GroupCreateSerizliser struct {
	Name    string   `json:"name" validate:"required,min=3,max=20"`
	Members []Member `json:"members" validate:"required,dive"`
}

type GroupGetSerizliser struct {
	IDParam
}

type AddMembersSerizliser struct {
	Members []Member `json:"members" validate:"required,dive"`
	IDParam
}

type removeMembersSerizliser struct {
	IDParam
	IDS []uint `json:"ids" validate:"required,dive"`
}