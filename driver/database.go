package driver

import (
	"database/sql"
	"io"
	"log"
	"os"

	// sql_driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//WriteLogFile : writing error to log file.
func WriteLogFile(err error) {
	f, erro := os.OpenFile("logs/output.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if erro != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err.Error())
		}
	}()

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	// log.Println(err)
}

// DbConn : Database connection
func DbConn() (db *sql.DB) {
	err := godotenv.Load()
	if err != nil {
		WriteLogFile(err)
		return
	}
	dbDriver := os.Getenv("DATABASE_DRIVER")
	dbUser := os.Getenv("DATABASE_USERNAME")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(:3306)/"+dbName)
	if err != nil {
		WriteLogFile(err)
		panic(err.Error())
	}
	return db
}
