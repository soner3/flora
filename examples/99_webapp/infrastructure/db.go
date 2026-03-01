package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/soner3/flora"
	"github.com/soner3/flora/examples/99_full_tutorial/config"

	_ "modernc.org/sqlite"
)

// DBConfig manages the database connection for the DI container.
type DBConfig struct {
	flora.Configuration
}

// ProvideDatabase connects to SQLite, creates a table, and inserts dummy data.
// flora:primary
func (c *DBConfig) ProvideDatabase(cfg *config.AppConfig) (*sql.DB, func(), error) {
	fmt.Printf("-> [DB] Connecting to database (%s)...\n", cfg.DSN)

	db, err := sql.Open("sqlite", cfg.DSN)
	if err != nil {
		return nil, nil, err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS books")
	if err != nil {
		return nil, nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY, title TEXT, author TEXT)`)
	if err != nil {
		return nil, nil, err
	}

	_, err = db.Exec(`INSERT INTO books (title, author) VALUES ('The Go Programming Language', 'Alan A. A. Donovan')`)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		fmt.Println("-> [DB] Closing database connection...")
		db.Close()
	}

	return db, cleanup, nil
}
