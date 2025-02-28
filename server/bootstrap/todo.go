package bootstrap

var todos []func()

func AddTodo(f func()) {
	todos = append(todos, f)
}

func InitTodo() {
	for _, todo := range todos {
		todo()
	}
}
