package dataModel

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq" //위의 sql에 대한 go-sqlite3를 사용한다는 암시적 의미로 추가
)

type pqHandler struct {
	db *sql.DB
}

func (s *pqHandler) GetTodos(sessionID string) []*Todo {
	todos := []*Todo{}
	// rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionID=?", sessionID)
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionID=$1", sessionID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}

	return todos
}
func (s *pqHandler) AddTodo(sessionID string, name string) *Todo {
	// stmt, err := s.db.Prepare("INSERT INTO todos (sessionID, name, completed, createdAt) VALUES (?, ?, ?, datetime('now'))")
	stmt, err := s.db.Prepare("INSERT INTO todos (sessionID, name, completed, createdAt) VALUES ($1, $2, $3, NOW()) RETURNING id")
	if err != nil {
		panic(err)
	}

	// result, err := stmt.Exec(sessionID, name, false)
	var id int
	err = stmt.QueryRow(sessionID, name, false).Scan(&id)
	if err != nil {
		panic(err)
	}

	// id, _ := result.LastInsertId() //마지막 추가된 아이디 값 나옴.
	var todo Todo
	// todo.ID = int(id)
	todo.ID = id
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()

	return &todo
}
func (s *pqHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=$1")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}

	cnt, _ := rst.RowsAffected()
	return cnt > 0
}
func (s *pqHandler) CompleteTodo(id int, complete bool) bool {
	stmt, err := s.db.Prepare("UPDATE todos SET completed=$1 WHERE id=$2")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(complete, id)
	if err != nil {
		panic(err)
	}
	cnt, _ := rst.RowsAffected()
	return cnt > 0
}
func (s *pqHandler) Close() {
	s.db.Close()
}
func newPQHandler(dbConn string) DataHandlerInterface {
	database, err := sql.Open("postgres", dbConn)
	if err != nil {
		panic(err)
	}

	statement, err := database.Prepare(
		// `CREATE TABLE IF NOT EXISTS todos (
		// id 			INTEGER PRIMARY KEY AUTOINCREMENT,
		// sessionID	STRING,
		// name		TEXT,
		// completed	BOOLEAN,
		// createdAt	DATETIME
		// );`)
		`CREATE TABLE IF NOT EXISTS todos (
			id 			SERIAL PRIMARY KEY,	
			sessionID	VARCHAR(256),
			name		TEXT,
			completed	BOOLEAN,
			createdAt	TIMESTAMP
		);`)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}

	// Index 작업. 기존 DB 순차탐색이라 느리기 떄문에 이렇게 처리하려 바이너리트리로 추가 될 때 Index 트리를 만듬.
	statement, err = database.Prepare(
		`CREATE INDEX IF NOT EXISTS sessionIDIndexOnTodos ON todos (
			sessionID ASC
		);`)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}

	return &pqHandler{db: database}
}
