package postgres

import (
	"database/sql"

	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
	_ "github.com/lib/pq"
)

// Establishes connection to postgres
func Connect(log *logger.Logger, url string) *sql.DB {
	log.Info(constants.PENDING + " Connecting to Postgres...")
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(constants.FAILURE + " Failed to connect to postgres: " + err.Error())
	}
	log.Info(constants.SUCCESS + " Connected to Postgres")
	return db
}
