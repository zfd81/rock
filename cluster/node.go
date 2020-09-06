package cluster

type Node struct {
	Id          string `json:"id"`
	Address     string `json:"addr"`
	Port        int    `json:"port"`
	StartUpTime int64  `json:"start-up-time"`
	LeaderFlag  bool   `json:"-"`
}

func NewNode(id string) *Node {
	return &Node{
		Id:         id,
		LeaderFlag: false,
	}
}
