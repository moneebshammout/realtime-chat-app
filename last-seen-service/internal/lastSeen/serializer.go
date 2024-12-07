package lastSeen

type IDParam struct {
    ID string `param:"id" validate:"required"`
}

type CreateSerizliser struct {
	UserId    any   `json:"userId" validate:"required,min=3"`
	SeenAt string `json:"seenAt" validate:"required"`
}