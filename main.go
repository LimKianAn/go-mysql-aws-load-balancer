package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db = openMysql()
	defer db.Close()
	printErr(db.Ping())

	http.HandleFunc("/", index)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/instanceID", instanceID)
	http.HandleFunc("/favicon.ico", http.NotFound)

	http.HandleFunc("/create", create)
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/read", read)
	http.HandleFunc("/drop", drop)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}

func index(responseWriter http.ResponseWriter, request *http.Request) {
	_, err := io.WriteString(responseWriter, "Hi there!")
	printErr(err)
}

func ping(responseWriter http.ResponseWriter, request *http.Request) {
	_, err := io.WriteString(responseWriter, "OK")
	printErr(err)
}

func instanceID(responseWriter http.ResponseWriter, request *http.Request) {
	response, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	printErr(err)

	bytes := make([]byte, response.ContentLength)
	response.Body.Read(bytes)
	response.Body.Close()
	io.WriteString(responseWriter, string(bytes))
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
