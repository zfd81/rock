package cluster

type Node struct {
	Id          string `json:"id"`
	Address     string `json:"addr"`
	Port        int    `json:"port"`
	StartUpTime int64  `json:"start-up-time"`
	LeaderFlag  bool   `json:"leader-flag"`
}

func NewNode(id string) *Node {
	return &Node{Id: id}
}
