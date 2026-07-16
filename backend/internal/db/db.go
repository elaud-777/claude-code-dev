package db

import (
	"strings"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	team_id INTEGER,
	created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS teams (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	invite_code TEXT UNIQUE NOT NULL,
	owner_id INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	team_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'TODO',
	creator_id INTEGER NOT NULL,
	assignee_id INTEGER,
	created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	team_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS ix_tasks_team_created ON tasks(team_id, created_at);
CREATE INDEX IF NOT EXISTS ix_messages_team_created ON messages(team_id, created_at);
CREATE INDEX IF NOT EXISTS ix_users_team_id ON users(team_id);
`

// postgresSchema uses SERIAL instead of AUTOINCREMENT for Postgres compatibility.
const postgresSchema = `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	team_id INTEGER,
	created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS teams (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	invite_code TEXT UNIQUE NOT NULL,
	owner_id INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
	id SERIAL PRIMARY KEY,
	team_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'TODO',
	creator_id INTEGER NOT NULL,
	assignee_id INTEGER,
	created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	team_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS ix_tasks_team_created ON tasks(team_id, created_at);
CREATE INDEX IF NOT EXISTS ix_messages_team_created ON messages(team_id, created_at);
CREATE INDEX IF NOT EXISTS ix_users_team_id ON users(team_id);
`

// Connect opens a DB connection based on databaseURL's scheme:
// "sqlite:///./path.db" (or a bare file path) uses the pure-Go SQLite driver,
// "postgres://..." uses the Postgres driver.
func Connect(databaseURL string) (*sqlx.DB, error) {
	if strings.HasPrefix(databaseURL, "postgres://") || strings.HasPrefix(databaseURL, "postgresql://") {
		conn, err := sqlx.Connect("postgres", databaseURL)
		if err != nil {
			return nil, err
		}
		if _, err := conn.Exec(postgresSchema); err != nil {
			return nil, err
		}
		return conn, nil
	}

	path := strings.TrimPrefix(databaseURL, "sqlite:///")
	path = strings.TrimPrefix(path, "sqlite://")
	if path == "" {
		path = "./taskflow.db"
	}
	conn, err := sqlx.Connect("sqlite", path)
	if err != nil {
		return nil, err
	}
	// SQLite allows only one writer at a time; serializing all connections
	// through a single pooled connection avoids spurious "database is locked"
	// errors under concurrent requests (e.g. Promise.all on the frontend).
	conn.SetMaxOpenConns(1)
	if _, err := conn.Exec(schema); err != nil {
		return nil, err
	}
	return conn, nil
}
