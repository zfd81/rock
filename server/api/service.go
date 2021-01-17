package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zfd81/rock/server"

	"github.com/zfd81/rock/server/services"

	"github.com/zfd81/rock/script"

	log "github.com/sirupsen/logrus"
	pb "github.com/zfd81/rock/proto/rockpb"

	"github.com/spf13/cast"
	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/meta/dai"
)

type Service struct{}

func (d *Service) TestAnalysis(ctx context.Context, request *pb.RpcRequest) (*pb.ServResponse, error) {
	name := request.Header["name"]
	source := request.Data
	serv, err := SourceAnalysis(source)
	if err != nil {
		return nil, fmt.Errorf("Bad parameter %s", err.Error())
	}
	serv.Name = name
	bytes, err := json.Marshal(serv)
	if err != nil {
		return nil, err
	}
	return &pb.ServResponse{
		Data: string(bytes),
	}, nil
}

func (d *Service) Test(ctx context.Context, request *pb.RpcRequest) (*pb.ServResponse, error) {
	source := request.Data
	serv, err := SourceAnalysis(source)
	if err != nil {
		return nil, fmt.Errorf("Bad syntax %s", err.Error())
	}
	serv.Name = request.Header["name"]
	serv.Source = source

	if serv.Method == "LOCAL" {

	}

	res := services.NewResource(serv).SetEnvironment(server.GetEnvironment())
	if len(request.Params) > 0 {
		for _, param := range res.GetParams() {
			val, found := request.Params[param.Name]
			if !found {
				return nil, fmt.Errorf("Parameter %s not found", param.Name)
			}
			if err = param.SetValue(val); err != nil {
				return nil, fmt.Errorf("Bad parameter %s data type error", param.Name)
			}
		}
	}
	log, resp, err := res.Run()
	if err != nil {
		return nil, fmt.Errorf(log)
	}
	bytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}
	return &pb.ServResponse{
		Message: log,
		Header:  resp.Header,
		Data:    string(bytes),
	}, nil
}

func (d *Service) CreateService(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	source := request.Data
	serv, err := SourceAnalysis(source)
	if err != nil {
		return nil, fmt.Errorf("Bad syntax %s", err.Error())
	}
	serv.Name = request.Params["name"]
	serv.Source = source
	err = dai.CreateService(serv)
	if err != nil {
		log.Error("Create service error: ", err)
		return nil, err
	}

	return &pb.RpcResponse{
		Code:    200,
		Message: fmt.Sprintf("Service %s created successfully", serv.Path),
	}, nil
}

func (d *Service) DeleteService(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	namespace := request.Params["namespace"]
	method := request.Params["method"]
	path := request.Params["path"]

	m := strings.ToUpper(method)
	if m != http.MethodGet && m != http.MethodPost &&
		m != http.MethodPut && m != http.MethodDelete &&
		m != "LOCAL" {
		return nil, fmt.Errorf("Method %s not found", method)
	}

	serv := &meta.Service{
		Namespace: namespace,
		Method:    method,
		Path:      path,
	}

	err := dai.DeleteService(serv)
	if err != nil {
		log.Error("Delete service error: ", err)
		return nil, err
	}

	return &pb.RpcResponse{
		Code:    200,
		Message: fmt.Sprintf("Service %s deleted successfully", serv.Path),
	}, nil
}

func (d *Service) ModifyService(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	source := request.Data
	serv, err := SourceAnalysis(source)
	if err != nil {
		return nil, fmt.Errorf("Bad syntax %s", err.Error())
	}
	serv.Source = source
	err = dai.ModifyService(serv)
	if err != nil {
		log.Error("Modify service error: ", err)
		return nil, err
	}

	return &pb.RpcResponse{
		Code:    200,
		Message: fmt.Sprintf("Service %s modified successfully", serv.Path),
	}, nil
}

func (d *Service) FindService(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	namespace := request.Params["namespace"]
	method := request.Params["method"]
	path := request.Params["path"]
	m := strings.ToUpper(method)
	if m != http.MethodGet && m != http.MethodPost &&
		m != http.MethodPut && m != http.MethodDelete &&
		m != "LOCAL" {
		return nil, fmt.Errorf("Method %s not found", method)
	}
	serv, err := dai.GetService(namespace, m, path)
	if err != nil {
		log.Error("Find service error: ", err)
		return nil, err
	}
	if serv != nil {
		serv.Source = ""
	}
	bytes, err := json.Marshal(serv)
	if err != nil {
		return nil, err
	}

	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}

func (d *Service) ListServices(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	namespace := request.Params["namespace"]
	path := request.Params["path"]
	servs, err := dai.ListService(namespace, path)
	if err != nil {
		log.Error("List services error: ", err)
		return nil, err
	}
	for _, serv := range servs {
		serv.Source = ""
	}
	bytes, err := json.Marshal(servs)
	if err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}

func SourceAnalysis(source string) (*meta.Service, error) {
	var definition string
	start := strings.Index(source, "$.define(")
	if start != -1 {
		end := strings.Index(source[start:], "})")
		if end == -1 {
			return nil, errs.New(errs.ErrParamBad, "Service definition error")
		}
		definition = source[start : end+3]
	}
	if definition == "" {
		return ModuleAnalysis(source)
	} else {
		serv := &meta.Service{}
		se := script.New()
		se.SetScript(definition)
		err := se.Run()
		if err != nil {
			return nil, errs.New(errs.ErrParamBad, "Service definition error:"+err.Error())
		}
		data, err := se.GetVar("__serv_definition")
		if err != nil {
			return nil, errs.New(errs.ErrParamBad, "Service definition error:"+err.Error())
		}
		define, ok := data.(map[string]interface{})
		if !ok {
			return nil, errs.New(errs.ErrParamBad, "Service definition error")
		}
		serv.Namespace = cast.ToString(define["namespace"])
		serv.Path = cast.ToString(define["path"])
		if serv.Path == "" {
			return nil, errs.New(errs.ErrParamBad, "Service path not found")
		}
		serv.Method = cast.ToString(define["method"])
		if serv.Method == "" {
			return nil, errs.New(errs.ErrParamBad, "Service method not found")
		}
		m := strings.ToUpper(serv.Method)
		if m != http.MethodGet && m != http.MethodPost &&
			m != http.MethodPut && m != http.MethodDelete {
			return nil, errs.New(errs.ErrParamBad, "Service method["+serv.Method+"] error")
		}
		params := define["params"]
		if params != nil {
			ps, ok := params.([]map[string]interface{})
			if !ok {
				return nil, errs.New(errs.ErrParamBad, "Service parameters definition error")
			}
			for _, param := range ps {
				serv.AddParam(cast.ToString(param["name"]), cast.ToString(param["dataType"]), cast.ToString(param["scope"]))
			}
		}
		return serv, nil
	}
}

func ModuleAnalysis(source string) (*meta.Service, error) {
	serv := &meta.Service{}
	m := services.NewModule(serv).SetEnvironment(server.GetEnvironment())
	se := script.NewWithContext(m)
	se.AddScript("var module = {};")
	se.AddScript(source)
	err := se.Run()
	if err != nil {
		return nil, errs.New(errs.ErrParamBad, "Module definition error:"+err.Error())
	}
	value, err := se.GetMlVar("module.exports.define")
	define, ok := value.(map[string]interface{})
	if !ok {
		return nil, errs.New(errs.ErrParamBad, "Module definition error")
	}
	serv.Namespace = cast.ToString(define["namespace"])
	serv.Path = cast.ToString(define["path"])
	if serv.Path == "" {
		return nil, errs.New(errs.ErrParamBad, "Module path not found")
	}
	serv.Method = "LOCAL"
	return serv, nil
}
