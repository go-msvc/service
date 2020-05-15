package http

import (
	"fmt"
	"net/http"

	"github.com/go-msvc/config"
	"github.com/go-msvc/service"
)

//server implements service.IServer
type server struct {
	config serverConstructor
}

func (s server) Serve(svc service.IService) error {
	if err := http.ListenAndServe(s.config.Address, Handler(svc)); err != nil {
		panic(err)
	}
	return nil
}

//serverConstructor implements service.IServerConstructor
type serverConstructor struct {
	Address string `json:"address"`
}

func (c *serverConstructor) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("missing address, e.g. \"localhost:12345\"")
	}
	return nil
}

func (c serverConstructor) Create() (service.IServer, error) {
	return server{
		config: c,
	}, nil
	//return nil, fmt.Errorf("NYI")
}

func init() {
	config.RegisterConstructor("http", &serverConstructor{}) //pass by reference so that Validate has ptr receivier
}
