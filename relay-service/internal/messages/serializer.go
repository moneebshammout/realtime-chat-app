package messages


type IGetUserMessages struct {
	ReceiverID     string `param:"id" validate:"required"`
}

type IDeleteMessages struct {
	IDS []string `query:"id" validate:"required,dive"`
}