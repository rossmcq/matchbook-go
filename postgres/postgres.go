package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = "localhost"
	port     = 5432
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = "matchbook"
)

func CheckConnection() error {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError("sql.Open", err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError("db.Ping", err)

	fmt.Println("DB Connected!")

	return nil
}

func CheckError(method string, err error) {
	if err != nil {
		fmt.Printf("Error with %v: %v", method, err)
	}
}
