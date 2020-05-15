package service

//IOper is implemented by user operations
type IOper interface {
	Handle(ctx IContext, req interface{}) (IResponse, error)
}

//IRequest ...
type IRequest interface{}

//IResponse ...
type IResponse interface{}

//IValidator ...
type IValidator interface {
	Validate() error
}
