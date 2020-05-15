package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	logger "github.com/go-msvc/log"
	"github.com/go-msvc/service"
)

var (
	log = logger.ForThisPackage()
)

//Handler implements http.Handler
func Handler(svc service.IService) http.Handler {
	return handler{svc: svc}
}

type handler struct {
	svc service.IService
}

func (h handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//prepare context
	ctx := service.Context(caller{})

	ctx.Debugf("HTTP %s %s\n", req.Method, req.URL.Path)
	operName := req.URL.Path[1:]
	ctx.Set("operName", operName)
	oper := h.svc.Oper(operName)
	if oper == nil {
		ctx.Debug("unknown operation")
		http.Error(res, "unknown oper("+operName+")", http.StatusNotFound)
		return
	}

	operReqPtrValue := reflect.New(reflect.TypeOf(oper))
	if err := json.NewDecoder(req.Body).Decode(operReqPtrValue.Interface()); err != nil && err != io.EOF {
		ctx.Debug("invalid JSON")
		http.Error(res, "invalid JSON", http.StatusNotFound)
		return
	}

	if validator, ok := operReqPtrValue.Interface().(service.IValidator); ok {
		if err := validator.Validate(); err != nil {
			ctx.Debug("invalid request")
			http.Error(res, "invalid request: "+err.Error(), http.StatusBadRequest)
			return
		}
		ctx.Debug("valid request")
	} else {
		ctx.Debug("no validator")
	}

	operRes, err := oper.Handle(ctx, operReqPtrValue.Elem().Interface())
	if err != nil {
		ctx.Debugf("oper failed: %v", err)
		http.Error(res, "failed", http.StatusInternalServerError)
		return
	}
	if operRes != nil {
		jsonRes, err := json.Marshal(operRes)
		if err != nil {
			ctx.Debugf("response marshalling failed: %v", err)
			http.Error(res, "failed", http.StatusNoContent)
			return
		}
		ctx.Debug("wrote response")
		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonRes)
	} else {
		ctx.Debug("no response")
	}
	ctx.Debug("done")
}

type caller struct {
}

func (c caller) Call(ctx service.IContext, name string, req interface{}) (service.IResponse, error) {
	url := "http://localhost:12345/" + name
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %v", url, err)
	}
	return res, nil
}
