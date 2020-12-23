package cmd

import (
	"log"
	"time"

	pb "github.com/zfd81/rock/proto/rockpb"

	"github.com/zfd81/rock/http"
	"google.golang.org/grpc"
)

type GlobalFlags struct {
	Endpoints []string
	User      string
	Password  string
}

var (
	client = http.New()
)

func url(path string) string {
	return globalFlags.Endpoints[0] + "/rock/" + path
}

func GetConnection() *grpc.ClientConn {
	address := globalFlags.Endpoints[0]
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func GetServiceClient() pb.ServiceClient {
	return pb.NewServiceClient(GetConnection())
}

func GetDataSourceClient() pb.DataSourceClient {
	return pb.NewDataSourceClient(GetConnection())
}

func GetClusterClient() pb.ClusterClient {
	return pb.NewClusterClient(GetConnection())
}
