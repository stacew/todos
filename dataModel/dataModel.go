package dataModel

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type DataHandlerInterface interface {
	GetTodos(sessionId string) []*Todo
	AddTodo(sessionId string, name string) *Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, complete bool) bool
	Close()
}

func NewDataHandler(filepath string) DataHandlerInterface {
	//return newMapHandler()
	return newSqliteHandler(filepath)
}
