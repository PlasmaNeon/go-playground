package main

import (
	"database/sql"
	"fmt"
	"geeorm"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func sqlTest() {
	db, _ := sql.Open("sqlite3", "../gee.db")
	defer func() { db.Close() }()
	db.Exec("DROP TABLE IF EXISTS User;")
	db.Exec("CREATE TABLE User(Name text);")
	result, err := db.Exec("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam")
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}
	row := db.QueryRow("SELECT Name FROM User limit 1")
	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	}
}

func main() {
	engine, _ := geeorm.NewEngine("sqlite3", "../gee.db")
	defer engine.Close()
	s := engine.NewSession()
	s.Raw("DROP TABLE IF EXISTS User;").Exec()
	s.Raw("CREATE TABLE User(Name text);").Exec()
	s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO USER(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
