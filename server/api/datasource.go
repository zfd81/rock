package api

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/meta/dai"
	pb "github.com/zfd81/rock/proto/rockpb"
)

type DataSource struct{}

func (d *DataSource) CreateDataSource(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	var ds meta.DataSource
	err := json.Unmarshal([]byte(request.Data), &ds)
	if err != nil {
		return nil, err
	}
	//数据源内容监测
	if ds.Name == "" ||
		ds.Driver == "" ||
		ds.Host == "" ||
		ds.Port < 100 ||
		ds.User == "" ||
		ds.Password == "" ||
		ds.Database == "" {
		return nil, fmt.Errorf("Bad parameter %s", "[DataSource information cannot be empty]")
	}
	if err = dai.CreateDataSource(&ds); err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code:    200,
		Message: fmt.Sprintf("DataSource %s created successfully", ds.Name),
	}, nil
}

func (d *DataSource) DeleteDataSource(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	namespace := request.Params["namespace"]
	name := request.Params["name"]
	ds := &meta.DataSource{
		Namespace: namespace,
		Name:      name,
	}
	if err := dai.DeleteDataSource(ds); err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code:    200,
		Message: fmt.Sprintf("DataSource %s deleted successfully", name),
	}, nil
}

func (d *DataSource) ModifyDataSource(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	var ds *meta.DataSource
	err := json.Unmarshal([]byte(request.Data), ds)
	if err != nil {
		return nil, err
	}
	//数据源内容监测
	if ds.Name == "" ||
		ds.Driver == "" ||
		ds.Host == "" ||
		ds.Port < 100 ||
		ds.User == "" ||
		ds.Password == "" ||
		ds.Database == "" {
		return nil, fmt.Errorf("Bad parameter %s", "[DataSource information cannot be empty]")
	}
	if err = dai.ModifyDataSource(ds); err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code:    200,
		Message: fmt.Sprintf("DataSource %s modified successfully", ds.Name),
	}, nil
}

func (d *DataSource) FindDataSource(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	namespace := request.Params["namespace"]
	name := request.Params["name"]
	ds, err := dai.GetDataSource(namespace, name)
	if err != nil {
		log.Error("Find datasource error: ", err)
		return nil, err
	}
	bytes, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}

func (d *DataSource) ListDataSources(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	namespace := request.Params["namespace"]
	dses, err := dai.ListDataSource(namespace, "/")
	if err != nil {
		log.Error("List datasources error: ", err)
		return nil, err
	}
	bytes, err := json.Marshal(dses)
	if err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}
