package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/rossmcq/matchbook-go/model"

	_ "github.com/lib/pq"
)

var (
	host     = "localhost"
	port     = 5432
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = "matchbook"
	psqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
)

func CheckConnection() error {
	// connection string

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("Error with sql.Open: %v", err)
	}

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("Error with db.Ping: %v", err)
	}

	fmt.Println("DB Connected!")

	return nil
}

func InsertOrReturnGameID(ctx context.Context, game model.Game) error {
	var gameID string

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("Error with sql.Open: %v", err)
	}

	// close database
	defer db.Close()

	selectStmt := `SELECT id FROM football.games 
					WHERE event_id = $1 
					AND market_id = $2;`

	row := db.QueryRow(selectStmt, game.EventID, game.MarketID)

	err = row.Scan(&gameID)
	fmt.Printf("Return GameID, row:%v: \n", gameID)
	if gameID != "" {
		fmt.Printf("gameID!=nilish: \n")

		return nil
	}

	insertDynStmt := `INSERT INTO football.games (id, event_id, market_id, description)
						VALUES ($1,$2,$3,$4);`
	_, err = db.Exec(insertDynStmt, game.GameID, game.EventID, game.MarketID, game.Description)
	if err != nil {
		return fmt.Errorf("Error with db.Exec: %v", err)
	}

	fmt.Printf("GAME INSERTED: %v \n", gameID)

	return nil
}
