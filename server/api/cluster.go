package api

import (
	"context"

	pb "github.com/zfd81/rock/proto/rockpb"
)

type ClusterServer struct{}

func (c *ClusterServer) ListMembers(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {

	//kvs, err := etcd.GetWithPrefix(cluster.MemberPath())
	//if err != nil {
	//	return nil, err
	//}
	//var nodes []*cluster.Node
	//var teams = make(map[string]*cluster.Team)
	//for _, kv := range kvs {
	//	n := &cluster.Node{}
	//	err := json.Unmarshal(kv.Value, n)
	//	if err != nil {
	//		return nil, err
	//	}
	//	var team *cluster.Team
	//	value, found := teams[n.Team]
	//	if found {
	//		team = value
	//	} else {
	//		team = &cluster.Team{}
	//		teams[n.Team] = team
	//	}
	//	team.AddMember(n)
	//	nodes = append(nodes, n)
	//}
	//for _, n := range nodes {
	//	team := teams[n.Team]
	//	n.LeaderFlag = team.IsLeader(n)
	//}
	//bytes, err := json.Marshal(nodes)
	//if err != nil {
	//	return nil, err
	//}
	return &pb.RpcResponse{
		Code: 200,
		//Data: string(bytes),
	}, nil
}

func (c *ClusterServer) MemberStatus(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	//tbls := server.GetDatabase("").Tables
	//tblInfos := []map[string]interface{}{}
	//for _, tbl := range tbls {
	//	info := map[string]interface{}{}
	//	info["name"] = tbl.Name
	//	cc, rc, size := tbl.Status()
	//	info["colCount"] = cc
	//	info["rowCount"] = rc
	//	info["tblSize"] = size
	//	tblInfos = append(tblInfos, info)
	//}
	//bytes, err := json.Marshal(tblInfos)
	//if err != nil {
	//	return nil, err
	//}
	return &pb.RpcResponse{
		Code: 200,
		//Data: string(bytes),
	}, nil
}
