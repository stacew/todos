package app

import (
	"net/http"
	"os"
	"stacew/todos/dataModel"

	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

var rd *render.Render = render.New()

type Success struct {
	Success bool `json:"success"`
}

type AppHandler struct {
	http.Handler //embeded is-a같은 has-a 관계라는데, 이름 정해주면 안 됨...
	dmHandler    dataModel.DataHandlerInterface
}

func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	list := a.dmHandler.GetTodos(sessionId)
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	name := r.FormValue("name")
	todo := a.dmHandler.AddTodo(sessionId, name)
	rd.JSON(w, http.StatusCreated, todo)
}

func (a *AppHandler) removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	ok := a.dmHandler.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func (a *AppHandler) completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	ok := a.dmHandler.CompleteTodo(id, complete)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func (a *AppHandler) Close() {
	a.dmHandler.Close()
}

//실제로 이 함수를 쓰긴 해야하지만, 테스트 함수에서 쓰기 위해 람다형식 개조.

// func getSessionID(r *http.Request) string {
var getSessionID = func(r *http.Request) string {
	session, _ := store.Get(r, "session")

	// Set some session values.
	val := session.Values["id"]
	if val == nil {
		return ""
	}

	return val.(string)
}

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//주의. user가 로그인 페이지를 눌렀는데, redirect signin.html로 무한루프 되기 때문에 이 부분을 처리해준다.
	//if request URL is /signin.html, then next() signin.css도 처리 필요
	if strings.Contains(r.URL.Path, "/signin") ||
		strings.Contains(r.URL.Path, "/auth") {
		next(w, r)
		return
	}

	sessionID := getSessionID(r)
	// if user already signed in
	if sessionID != "" {
		next(w, r) //다음 핸들러 static file
		return
	}

	// if not user sign in redirect signin.html
	http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)
}

func MakeNewHandler(dbConn string) *AppHandler {

	r := mux.NewRouter()

	//네그로니에 미들웨어를 추가할 수 있다. 구글 세션 때문에 데코레이터 패턴을 새로 안 만들고 사용 가능.
	//네그로니 Classic()에 3가지 데코레이터 기본 recovery, logger, static file 처리가 있는데
	//staic file 처리 전에 Signin 체크를 먼저 해야하기 때문에 밖으로 New를 뺌.
	//neg := negroni.Classic()
	//siginin에서 끊기면 router로 돌아가도록 negroni.HandlerFunc(CheckSignin) 네그로니의 핸들러 펑션으로 동록
	neg := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.HandlerFunc(CheckSignin), negroni.NewStatic(http.Dir("public")))

	neg.UseHandler(r)

	a := &AppHandler{
		Handler:   neg,
		dmHandler: dataModel.NewDataHandler(dbConn),
	}

	r.HandleFunc("/todoH", a.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todoH", a.addTodoHandler).Methods("POST")
	r.HandleFunc("/todoH/{id:[0-9]+}", a.removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/complete-todoH/{id:[0-9]+}", a.completeTodoHandler).Methods("GET")
	r.HandleFunc("/auth/google/login", googleLoginHandler)
	r.HandleFunc("/auth/google/callback", googleAuthCallback)
	r.HandleFunc("/", a.indexHandler)

	return a
}
