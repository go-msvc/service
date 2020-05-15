package service

import (
	"fmt"
)

//New ...
func New() IService {
	return service{
		operByName: map[string]IOper{},
	}
}

//IService ...
type IService interface {
	WithOper(name string, oper IOper) IService
	Oper(name string) IOper
	Run() error
}

type service struct {
	operByName map[string]IOper
}

func (s service) WithOper(name string, oper IOper) IService {
	if _, ok := s.operByName[name]; ok {
		panic(fmt.Errorf("duplicate oper name \"%s\"", name))
	}
	s.operByName[name] = oper
	return s
}

func (s service) Oper(name string) IOper {
	if so, ok := s.operByName[name]; ok {
		return so
	}
	return nil
}

func (s service) Run() error {
	return fmt.Errorf("no server to run")
}
