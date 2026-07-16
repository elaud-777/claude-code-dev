package app

import (
	"github.com/jmoiron/sqlx"

	"taskflow-backend/internal/config"
)

// App holds shared dependencies (DB connection, settings) used by all handlers.
type App struct {
	DB       *sqlx.DB
	Settings config.Settings
}

// Q rebinds a query written with `?` placeholders to the bind style required
// by the currently connected driver (SQLite keeps `?`, Postgres becomes `$1`, ...).
func (a *App) Q(query string) string {
	return a.DB.Rebind(query)
}
