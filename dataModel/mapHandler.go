package dataModel

import "time"

type mapHandler struct {
	todoMap map[int]*Todo
	genID   int
}

//var todoSlice []*Todo

func (m *mapHandler) GetTodos(sessionId string) []*Todo {
	list := []*Todo{}
	for _, v := range m.todoMap {
		list = append(list, v)
	}
	return list
}
func (m *mapHandler) AddTodo(sessionId string, name string) *Todo {
	id := m.genID
	m.genID++

	todo := &Todo{id, name, false, time.Now()}
	m.todoMap[id] = todo
	return todo
}
func (m *mapHandler) RemoveTodo(id int) bool {
	if _, ok := m.todoMap[id]; ok {
		delete(m.todoMap, id)
		return true
	}
	return false
}
func (m *mapHandler) CompleteTodo(id int, complete bool) bool {

	if todo, ok := m.todoMap[id]; ok {
		todo.Completed = complete
		return true
	}
	return false
}
func (m *mapHandler) Close() {
}

func newMapHandler() DataHandlerInterface {
	m := &mapHandler{}
	m.todoMap = make(map[int]*Todo)
	m.genID = 0
	return m
}
