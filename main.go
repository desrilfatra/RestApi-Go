package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"restapi-go/handler"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "RestAPI-Go"
)

var (
	db *sql.DB

	err error
)

const PORT = ":8080"

func main() {
	db, err = sql.Open("postgres", ConnectDbPsql(host, user, password, dbname, port))
	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	r := mux.NewRouter()
	// r.HandleFunc("/", HomeHandler)
	itemHandler := handler.NewItemHandler(db)
	r.HandleFunc("/order", itemHandler.ItemHandler)
	r.HandleFunc("/order/{item_id}", itemHandler.ItemHandler)

	fmt.Println("Server is running 0.0.0.0" + PORT)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func ConnectDbPsql(host string, user string, password string, dbname string, port int) string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname)
	return psqlInfo
}
