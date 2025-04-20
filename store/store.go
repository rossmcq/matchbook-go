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

// DbConnection represents a connection to the PostgreSQL database
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

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.host, dbCfg.port, dbCfg.user, dbCfg.password, dbCfg.dbname)

	dbConnection, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return DbConnection{}, fmt.Errorf("failed to open database connection: %w", err)
	}

	return DbConnection{
		Database: dbConnection,
	}, nil
}

// CheckConnection verifies that the database connection is alive and working
func (d DbConnection) CheckConnection() error {
	if err := d.Database.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Database connection verified")
	return nil
}

func (d DbConnection) SelectGame(ctx context.Context, game model.Game) (string, error) {
	var gameID string

	selectStmt := `SELECT id FROM football.games 
					WHERE event_id = $1 
					AND market_id = $2;`
	row := d.Database.QueryRowContext(ctx, selectStmt, game.EventID, game.MarketID)

	err := row.Scan(&gameID)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", fmt.Errorf("failed to scan game ID: %w", err)
		}
	}

	return gameID, nil
}

func (d DbConnection) GetOpenGames(ctx context.Context) ([]model.Game, error) {
	var games []model.Game

	selectStmt := `SELECT id, event_id, market_id, start_at, status, home_team, away_team, description 
					FROM football.games WHERE status = 'open';`
	rows, err := d.Database.QueryContext(ctx, selectStmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query open games: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var game model.Game
		if err := rows.Scan(&game.GameID, &game.EventID, &game.MarketID, &game.StartAt, &game.Status, &game.HomeTeam, &game.AwayTeam, &game.Description); err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}
		games = append(games, game)
	}

	return games, nil
}

// CreateGame creates a new game in the database if it doesn't already exist
func (d DbConnection) CreateGame(ctx context.Context, game model.Game) error {
	gameID, err := d.SelectGame(ctx, game)
	if err != nil {
		return fmt.Errorf("failed to select game: %w", err)
	}
	if gameID != "" {
		log.Printf("Game already exists with ID: %s", gameID)
		return nil
	}

	insertDynStmt := `INSERT INTO football.games (id, event_id, market_id, start_at, status, home_team, away_team, description)
						VALUES ($1,$2,$3,$4,$5,$6,$7,$8);`

	_, err = d.Database.ExecContext(ctx, insertDynStmt,
		game.GameID, game.EventID, game.MarketID, game.StartAt,
		game.Status, game.HomeTeam, game.AwayTeam, game.Description)

	if err != nil {
		return fmt.Errorf("failed to insert game: %w", err)
	}

	log.Printf("Successfully created game with ID: %s", game.GameID)
	return nil
}
