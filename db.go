package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	db       *sql.DB
	fileName string
	filePath string
}

func (d Db) getDb() *sql.DB {
	return d.db
}

func (d *Db) initDb() {
	var err error
	rawPath := d.filePath + d.fileName
	path := filepath.Clean(rawPath)
	fmt.Println(path)
	(*d).db, err = sql.Open("sqlite3", path)
	fmt.Printf("%p\n", (*d).db)
	if err != nil {
		log.Fatal("not able to open database")
	}
}

func (d Db) closeDb() {
	d.db.Close()
	fmt.Println("closing DB")
}

func InitDb(fN, fP string) *Db {
	retVal := Db{nil, fN, fP}
	retVal.initDb()
	return &retVal
}

func AddTable(d *Db, tableName string, columns []string) {
	fmt.Printf("db ptr: %p \n", d)
	var columnString string

	for i, column := range columns {

		var primaryKey bool = false
		var endKey bool = false
		if i == 0 {
			primaryKey = true
		}
		if i == len(columns)-1 {
			endKey = true
		}

		values := strings.Split(column, " ")
		if len(values) != 2 {
			log.Fatal("improper schema")
		}
		columnString = columnString + values[0] + " " + values[1] + " NOT NULL "
		if primaryKey {
			columnString += " PRIMARY KEY "
		}

		if endKey == false {
			columnString += ","
		}

		columnString += "\n"
	}

	query := "CREATE TABLE IF NOT EXISTS " + tableName + " ( " + columnString + " );"
	fmt.Printf("%p\n", (*d).db)
	fmt.Println(query)
	_, err := d.getDb().Exec(query)
	if err != nil {
		log.Fatal("unable to add table")
	}
}

func AddNewPassword(d *Db, tableName, username, hash string) {

	// SQL query to insert username and hash
	query := `INSERT INTO ` + tableName + `(username, hash) VALUES (?, ?);`

	// Execute the query with the username and hashed password
	_, err := d.getDb().Exec(query, username, hash)
	if err != nil {
		log.Fatal("error inserting user: %w", err)
	}
}

func GetUserNamePassword(d *Db, tableName, username string) string {

	var hash string

	query := "SELECT hash FROM " + tableName + " where username = '" + username + "'"
	fmt.Println(query)
	rows, err := d.getDb().Query(query)
	if err != nil {
		log.Fatal("error while querying db")
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&hash); err != nil {
			log.Fatal(err)
		}
	}
	return hash
}

/*
func initPassDb() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS password_table (
		username TEXT NOT NULL PRIMARY KEY,
		hash TEXT NOT NULL
	);`)

	query := `INSERT INTO password_table (username, hash) VALUES (shivam, password);`

	db.Exec(query)
	fmt.Println("created table")

	if err != nil {
		log.Fatal("unable to create db")
	}
}
*/
