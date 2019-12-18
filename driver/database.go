package driver
import(
      "database/sql"
      "os"

    _ "github.com/go-sql-driver/mysql"
)



func DbConn() (db *sql.DB) { //Database connection
	dbDriver := os.Getenv("DATABASE_DRIVER")
	dbUser := os.Getenv("DATABASE_USERNAME")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
