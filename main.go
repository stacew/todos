package main

import (
	"net/http"

	"github.com/stacew/todos/app"
)

func main() {

	mux := app.MakeNewHandler("./test.db") //실행인자 이용 가능 flag.Arg
	defer mux.Close()

	http.ListenAndServe(":8080", mux)
}
