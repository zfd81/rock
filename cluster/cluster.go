package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/zfd81/parrot/conf"
	"github.com/zfd81/parrot/util/etcd"
)

var (
	leaseID  clientv3.LeaseID
	node     *Node
	leaderId string
	members  = make(map[string]*Node)
	config   = conf.GetConfig()
)

func GetNode() *Node {
	return node
}

func GetMembers() map[string]*Node {
	return members
}

func Register(startUpTime int64) error {
	ip, err := externalIP()
	if err != nil {
		return err
	}
	cli := etcd.GetClient()
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(config.Cluster.HeartbeatInterval))
	if err != nil {
		return err
	}
	node = NewNode(fmt.Sprintf("%x", session.Lease()))
	node.Address = ip.String()
	node.Port = config.Port
	node.StartUpTime = startUpTime

	//获得集群根目录
	path := config.Cluster.LeaderPath

	//获得集群成员结点目录
	mpath := config.Cluster.MemberPath

	//监听集群leader结点变化
	etcd.Watch(path, func(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
		leaderId = string(value)
	})

	//监听集群结点变化
	etcd.WatchWithPrefix(mpath, clusterWatcher)

	//加载现有结点
	kvs, err := etcd.GetWithPrefix(mpath)
	if err == nil {
		for _, kv := range kvs {
			n := addNode(kv.Key, kv.Value)
			if n.LeaderFlag {
				leaderId = n.Id
			}
		}
	}

	//结点注册并参与选举
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}
	elect := concurrency.NewElection(session, mpath)
	go func() {
		//竞选 Leader，直到成为 Leader 函数才返回
		if err := elect.Campaign(context.Background(), string(data)); err != nil {
			fmt.Println(err)
		} else {
			node.LeaderFlag = true
			if _, err = etcd.PutWithLease(path, node.Id, session.Lease()); err != nil {
				fmt.Println(err)
			}
		}
	}()
	return err
}

func clusterWatcher(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
	if operType == etcd.CREATE {
		addNode(key, value)
	} else if operType == etcd.MODIFY {
	} else if operType == etcd.DELETE {
		delete(members, string(key))
	}
}

func addNode(key []byte, value []byte) *Node {
	node := &Node{}
	err := json.Unmarshal(value, node)
	if err == nil {
		members[node.Id] = node
	}
	return node
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}
	return ip
}
