package messages


type IGetUserMessages struct {
	ReceiverID     string `param:"id" validate:"required"`
}