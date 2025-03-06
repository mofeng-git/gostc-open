package todo

import (
	"server/repository"
	"server/service/common/node_port"
)

func gostNodePort() {
	db, _, _ := repository.Get("")
	go node_port.Run(db)
}
