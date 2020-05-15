package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-msvc/log"
)

var (
	idGen  IIDGenerator
	ctxLog = log.Logger("ctx")
)

func init() {
	idGen = &defaultGen{}
}

//Context starts a blank new context
func Context(c ICaller) IContext {
	//assign a unique uuid
	id := idGen.New()
	return msContext{
		Context: context.Background(),
		ILogger: ctxLog.Temp(id),
		caller:  c,
		id:      id,
	}
}

//IContext ...
type IContext interface {
	context.Context
	log.ILogger
	ID() string
	Call(name string, req interface{}) (IResponse, error)
}

//msContext implements IContext
type msContext struct {
	context.Context
	log.ILogger
	caller ICaller
	id     string
}

func (ctx msContext) ID() string { return ctx.id }

func (ctx msContext) Call(name string, req interface{}) (IResponse, error) {
	if ctx.caller == nil {
		return nil, fmt.Errorf("no caller installed in this context")
	}
	res, err := ctx.caller.Call(ctx, name, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call %s: %v", name, err)
	}
	return res, nil
}

//IIDGenerator is used to generate unique ids
type IIDGenerator interface {
	New() string
}

type defaultGen struct {
	mutex sync.Mutex
	last  int
}

func (d *defaultGen) New() string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.last++
	return fmt.Sprintf("%08x", d.last)
}

//ICaller makes service calls
type ICaller interface {
	Call(ctx IContext, name string, req interface{}) (IResponse, error)
}
