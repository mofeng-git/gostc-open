package todo

import (
	"server/bootstrap"
	"server/repository"
	"server/service/common/node_port"
)

func init() {
	bootstrap.AddTodo(func() {
		db, _, _ := repository.Get("")
		go node_port.Run(db)
	})
}
