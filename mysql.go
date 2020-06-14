package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var db *sql.DB

func openMysql() *sql.DB {
	godotenv.Load(".env")
	password := os.Getenv("AWS_MYSQL_INSTANCE_PASSWORD")
	domainName := "go-mysql-instance.c25xrnjxvftp.eu-central-1.rds.amazonaws.com"
	dbName := "go-mysql"

	dataSource := "admin:" + password + "@tcp(" + domainName + ":3306)/"
	db, err := sql.Open("mysql", dataSource)
	printErr(err)

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	printErr(err)

	_, err = db.Exec("USE " + dbName)
	printErr(err)

	return db
}

func create(responseWriter http.ResponseWriter, request *http.Request) {
	table := request.FormValue("table")
	column := request.FormValue("column")
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS " + table + " (" + column + " VARCHAR(20))")
	printErr(err)
	defer statement.Close()

	result, err := statement.Exec()
	printErr(err)

	n, err := result.RowsAffected()
	printErr(err)

	fmt.Fprintln(responseWriter, "New table created.", n, "row(s) affected.")
}

func insert(responseWriter http.ResponseWriter, request *http.Request) {
	table := request.FormValue("table")
	column := request.FormValue("column")
	value := request.FormValue("value")
	str := "INSERT INTO " + table + "(" + column + ") VALUES (\"" + value + "\")"
	statement, err := db.Prepare(str)
	printErr(err)
	defer statement.Close()

	result, err := statement.Exec()
	printErr(err)

	n, err := result.RowsAffected()
	printErr(err)

	fmt.Fprintln(responseWriter, "New value inserted.", n, "row(s) affected.")
}

func read(responseWriter http.ResponseWriter, request *http.Request) {
	table := request.FormValue("table")
	column := request.FormValue("column")
	rows, err := db.Query("SELECT " + column + " FROM " + table)
	printErr(err)
	defer rows.Close()

	for rows.Next() {
		row := ""
		err := rows.Scan(&row)
		printErr(err)
		fmt.Fprintln(responseWriter, row)
	}
}

func update(responseWriter http.ResponseWriter, request *http.Request) {
	table := request.FormValue("table")
	column := request.FormValue("column")
	newValue := request.FormValue("new-value")
	oldValue := request.FormValue("old-value")
	str := "UPDATE " + table + " SET " + column + "=\"" + newValue + "\" WHERE " + column + "=\"" + oldValue + "\""
	statement, err := db.Prepare(str)
	printErr(err)
	defer statement.Close()

	results, err := statement.Exec()
	printErr(err)

	n, err := results.RowsAffected()
	printErr(err)

	fmt.Fprintln(responseWriter, n, "row(s) updated.")
}

func delete(responseWriter http.ResponseWriter, request *http.Request) {
	table := request.FormValue("table")
	column := request.FormValue("column")
	value := request.FormValue("value")
	statement, err := db.Prepare(`DELETE FROM ` + table + ` WHERE ` + column + `="` + value + `"`)
	printErr(err)
	defer statement.Close()

	results, err := statement.Exec()
	printErr(err)

	n, err := results.RowsAffected()
	printErr(err)

	fmt.Fprintln(responseWriter, n, "row(s) deleted.")
}

func drop(responseWriter http.ResponseWriter, request *http.Request) {
	table := request.FormValue("table")
	statement, err := db.Prepare(`DROP TABLE ` + table)
	printErr(err)
	defer statement.Close()

	results, err := statement.Exec()
	printErr(err)

	n, err := results.RowsAffected()
	printErr(err)

	fmt.Fprintln(responseWriter, table, "dropped.", n, "row(s) affected.")
}
