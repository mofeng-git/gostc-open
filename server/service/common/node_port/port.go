package node_port

import (
	"errors"
	"server/repository/query"
	"sync"
	"time"
)

var ports = make(map[string]map[string]bool)
var lock = &sync.Mutex{}

func Run(db *query.Query) {
	for {
		arrange(db)
		time.Sleep(time.Hour * 2)
	}
}

func GetPort(nodeCode string) (string, error) {
	lock.Lock()
	defer lock.Unlock()
	nodePorts, ok := ports[nodeCode]
	if !ok {
		return "", errors.New("端口资源不足")
	}
	if len(nodePorts) == 0 {
		return "", errors.New("端口资源不足")
	}
	for port, _ := range nodePorts {
		delete(ports[nodeCode], port)
		return port, nil
	}
	return "", errors.New("端口资源不足")
}

func ValidPort(nodeCode, port string, remove bool) bool {
	lock.Lock()
	defer lock.Unlock()
	var flag = ports[nodeCode][port]
	if flag && remove {
		delete(ports[nodeCode], port)
	}
	return flag
}

func ReleasePort(nodeCode string, port string) {
	lock.Lock()
	defer lock.Unlock()
	ports[nodeCode][port] = true
}

func arrange(db *query.Query) {
	lock.Lock()
	defer lock.Unlock()
	nodes, _ := db.GostNode.Find()
	nodePorts, _ := db.GostNodePort.Find()
	var nodeUsedPorts = make(map[string][]string)
	for _, nodePort := range nodePorts {
		nodeUsedPorts[nodePort.NodeCode] = append(nodeUsedPorts[nodePort.NodeCode], nodePort.Port)
	}
	for _, node := range nodes {
		ports[node.Code] = node.GetPorts(nodeUsedPorts[node.Code])
	}
}

func Arrange(db *query.Query, code string) {
	lock.Lock()
	defer lock.Unlock()
	node, _ := db.GostNode.Where(db.GostNode.Code.Eq(code)).First()
	var nodeUsedPorts []string
	_ = db.GostNodePort.Where(db.GostNodePort.NodeCode.Eq(code)).Pluck(db.GostNodePort.Port, &nodeUsedPorts)
	ports[node.Code] = node.GetPorts(nodeUsedPorts)
}
