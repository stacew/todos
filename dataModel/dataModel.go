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

func NewDataHandler(dbConn string) DataHandlerInterface {
	return newMapHandler()
	//return newSqliteHandler(filepath)
	//return newPQHandler(dbConn) //file path 에서 디비 접속 정보가 와야 함.
}
