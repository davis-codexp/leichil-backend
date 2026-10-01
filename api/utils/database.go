package utils

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

var db *sql.DB

func GetDBInstance() *sql.DB {
	return db
}

func GetConnection() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Error("Unable to Load Database Credentials from env: ", err.Error())
		return nil, err
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("PASSWORD")
	database := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=false", user, password, host, port, database)
	DB, err := sql.Open("mysql", uri)
	db = DB
	if err != nil {
		fmt.Println("Unable to connect to MySQL", err.Error())
		return nil, err
	}
	log.Info("Connected to MySQL")
	return db, nil
}

func RunQuery[T any](query string, fields []any, item *T, args []any) ([]any, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Error("Error in DB Query: ", query, ": ", err.Error())
		return nil, err
	}
	var result []any
	for rows.Next() {
		err = rows.Scan(fields...)
		if err != nil {
			log.Error("Error in scanning Rows: ", err.Error())
			return nil, err
		}
		result = append(result, *item)
	}
	return result, nil
}

func RunInsert(query string, fields []any) bool {
	var flag = true
	insert, err := db.Prepare(query)
	if err != nil {
		log.Error("Error in preparing Insert Query: ", err.Error())
		return false
	}
	defer insert.Close()
	_, err = insert.Exec(fields...)
	if err != nil {
		log.Error("Error in Insert Query Execution: ", err.Error())
		flag = false
	}
	return flag
}

func RunUpdate(query string, fields []any) bool {
	_, err := db.Exec(query, fields...)
	if err != nil {
		log.Error("Error in Update Query Execution: ", err.Error())
		return false
	}
	return true
}

func RunDelete(query string, args []any) error {
	_, err := db.Query(query, args...)
	return err
}
