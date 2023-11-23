package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rossmcq/matchbook-go/model"

	_ "github.com/lib/pq"
)

type DbConnection struct {
	Database *sql.DB
}


func (d *DbConnection) CheckConnection() error {
	// check db
	err := d.Database.Ping()

	if err != nil {
		return fmt.Errorf("Error with db.Ping: %v", err)
	}

	fmt.Println("DB Connected!")

	return nil
}

func (d DbConnection) InsertOrReturnGameID(ctx context.Context, game model.Game) error {
	var gameID string

	selectStmt := `SELECT id FROM football.games 
					WHERE event_id = $1 
					AND market_id = $2;`

	row := d.Database.QueryRow(selectStmt, game.EventID, game.MarketID)

	err := row.Scan(&gameID)
	fmt.Printf("Return GameID, row:%v: \n", gameID)
	if gameID != "" {
		fmt.Printf("gameID!=nilish: \n")

		return nil
	}

	insertDynStmt := `INSERT INTO football.games (id, event_id, market_id, description)
						VALUES ($1,$2,$3,$4);`
	_, err = d.Database.Exec(insertDynStmt, string(game.GameID), game.EventID, game.MarketID, game.Description)
	if err != nil {
		return fmt.Errorf("Error with db.Exec: %v", err)
	}


	fmt.Printf("GAME INSERTED: %v \n", gameID)

	return nil
}
