package main

import (
	"log"
	"net/http"
	"os"
	"stacew/todos/app"
)

func main() {

	//mux := app.MakeNewHandler("./test.db") //실행인자 이용 가능 flag.Arg
	mux := app.MakeNewHandler(os.Getenv("DATABASE_URL"))
	defer mux.Close()

	log.Println("Started App")
	err := http.ListenAndServe(":"+GetPort(), mux)
	if err != nil {
		panic(err)
	}

}

//GetPort is
func GetPort() string {
	port := os.Getenv("PORT") //env port
	if port == "" {
		port = "8080" //local
	}
	log.Print("port:" + port)
	return port
}
