package node_port

import (
	"errors"
	"gorm.io/gorm"
	"server/model"
	"sync"
	"time"
)

var ports = make(map[string][]string)
var lock = &sync.Mutex{}

func Run(db *gorm.DB) {
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
	var newPorts = make([]string, len(nodePorts)-1)
	copy(newPorts, nodePorts[1:])
	ports[nodeCode] = newPorts
	return nodePorts[0], nil
}

func ReleasePort(nodeCode string, port string) {
	lock.Lock()
	defer lock.Unlock()
	ports[nodeCode] = append(ports[nodeCode], port)
}

func arrange(db *gorm.DB) {
	lock.Lock()
	defer lock.Unlock()
	var nodes []model.GostNode
	db.Find(&nodes)
	var nodePorts []model.GostNodePort
	db.Find(&nodePorts)
	var nodeUsedPorts = make(map[string][]string)
	for _, nodePort := range nodePorts {
		nodeUsedPorts[nodePort.NodeCode] = append(nodeUsedPorts[nodePort.NodeCode], nodePort.Port)
	}
	for _, node := range nodes {
		ports[node.Code] = node.GetPorts(nodeUsedPorts[node.Code])
	}
}

func Arrange(db *gorm.DB, code string) {
	lock.Lock()
	defer lock.Unlock()
	var node model.GostNode
	db.Where("code = ?", code).First(&node)
	var nodeUsedPorts []string
	db.Model(&model.GostNodePort{}).Where("node_code = ?", code).Pluck("port", &nodeUsedPorts)
	ports[node.Code] = node.GetPorts(nodeUsedPorts)
}
