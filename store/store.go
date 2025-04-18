package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

func (d DbConnection) CheckConnection() error {
	// check db
	err := d.Database.Ping()

	if err != nil {
		return fmt.Errorf("error with db.Ping: %v", err)
	}

	log.Println("DB Connected!")

	return nil
}

func (d DbConnection) CreateGame(ctx context.Context, game model.Game) error {
	var gameID string

	selectStmt := `SELECT id FROM football.games 
					WHERE event_id = $1 
					AND market_id = $2;`

	row := d.Database.QueryRow(selectStmt, game.EventID, game.MarketID)

	err := row.Scan(&gameID)
	if err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("error scanning gameID: %s", err)
		}
	}
	fmt.Printf("returned GameID from db: %v: \n", gameID)
	if gameID != "" {
		fmt.Printf("gameID!=nilish: \n")

		return nil
	}

	insertDynStmt := `INSERT INTO football.games (id, event_id, market_id, start_at, status, home_team, away_team, description)
						VALUES ($1,$2,$3,$4,$5,$6,$7,$8);`
	_, err = d.Database.Exec(insertDynStmt, game.GameID, game.EventID, game.MarketID, game.StartAt, game.Status, game.HomeTeam, game.AwayTeam, game.Description)
	if err != nil {
		return fmt.Errorf("error inserting row to football.games: %v", err)
	}

	log.Printf("game inserted: %v \n", gameID)

	return nil
}
