package driver
import(
      "database/sql"

    _ "github.com/go-sql-driver/mysql"
)



func DbConn() (db *sql.DB) { //Database connection
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "Root@987"
	dbName := "weekly_update"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
