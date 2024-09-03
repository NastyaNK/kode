package connection

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"myproject/pkg/config"
)

type Connection struct {
	*sqlx.DB
}

func NewConnection(database *config.Database) (*Connection, error) {
	db, err := sqlx.Connect("postgres", database.GetDataSource())
	if err != nil {
		slog.Error("Failed to connect to the database",
			"error", err,
			"db_name", database.Database,
		)
		return nil, err
	}
	slog.Info("Connected to the database", "db_name", database.Database)
	return &Connection{db}, nil
}
