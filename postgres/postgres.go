package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/rossmcq/matchbook-go/model"

	_ "github.com/lib/pq"
)

type DbConnection struct {
	Database *sql.DB
}

type dbConfig struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func New() (DbConnection, error) {
	dbCfg := dbConfig{
		host:     "localhost",
		port:     5432,
		user:     "postgres",
		password: "postgres",
		dbname:   "matchbook",
	}

	if user, exists := os.LookupEnv("POSTGRES_USER"); exists {
		dbCfg.user = user
	}

	if password, exists := os.LookupEnv("POSTGRES_PASSWORD"); exists {
		dbCfg.password = password
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbCfg.host, dbCfg.port, dbCfg.user, dbCfg.password, dbCfg.dbname)

	dbConnection, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return DbConnection{}, fmt.Errorf("can't open DB: %v", err)
	}

	return DbConnection{
		Database: dbConnection,
	}, nil
}

func (d *DbConnection) CheckConnection() error {
	// check db
	err := d.Database.Ping()

	if err != nil {
		return fmt.Errorf("error with db.Ping: %v", err)
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
	if err != nil {
		return fmt.Errorf("error scanning gameID: %s", err)
	}
	fmt.Printf("returned GameID from DB: %v: \n", gameID)
	if gameID != "" {
		fmt.Printf("gameID!=nilish: \n")

		return nil
	}

	insertDynStmt := `INSERT INTO football.games (id, event_id, market_id, description)
						VALUES ($1,$2,$3,$4);`
	_, err = d.Database.Exec(insertDynStmt, fmt.Sprint(game.GameID), game.EventID, game.MarketID, game.Description)
	if err != nil {
		return fmt.Errorf("error with db.Exec: %v", err)
	}

	fmt.Printf("GAME INSERTED: %v \n", gameID)

	return nil
}
