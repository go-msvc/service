package service

import (
	"fmt"

	"github.com/go-msvc/config"
	"github.com/go-msvc/config/source/static"
	"github.com/pkg/errors"
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
	var server IServer
	var err error

	//ad default static config for an HTTP if no other server is configured
	config.Sources().Add(static.New("service/server/http", defaultConfig{Address: "localhost:8080"}))
	err = config.Create("service/server", &server)
	if err != nil {
		return errors.Wrapf(err, "failed to construct server")
	}
	if err := server.Serve(s); err != nil {
		return errors.Wrapf(err, "server failed")
	}
	return nil
}

type serviceConfig struct {
	Server IServerConstructor `json:"server"`
}

func (c serviceConfig) Validate() error {
	return nil
}

//IServerConstructor ...
type IServerConstructor interface {
	Create() (IServer, error)
}

//IServer ...
type IServer interface {
	config.IConstructed
	Serve(IService) error
}

type defaultConfig struct {
	Address string `json:"address"`
}

func (c defaultConfig) Validate() error { return nil }
