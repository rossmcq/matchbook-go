package application

import (
	"fmt"
	"os"

	"github.com/rossmcq/matchbook-go/matchbook"
)

type Config struct {
	dbConnectionString string
	matchbookToken     string
}

type dbConfig struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func LoadConfig() Config {
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

	matchbookToken, err := matchbook.LoadMatchboookToken()
	if err != nil {
		fmt.Printf("Error loading matchbook token %s", err)
	}
	return Config{
		dbConnectionString: psqlconn,
		matchbookToken:     *matchbookToken}
}
