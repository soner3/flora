// Copyright © 2026 Soner Astan astansoner@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/soner3/flora"
	"github.com/soner3/flora/examples/99_webapp/config"

	_ "modernc.org/sqlite"
)

// DBConfig manages the database connections for the DI container.
type DBConfig struct {
	flora.Configuration
}

// ProvideMasterDB provides the Master database (intended for write operations).
// flora:name=masterDB
func (c *DBConfig) ProvideMasterDB(cfg *config.AppConfig) (*sql.DB, func(), error) {
	fmt.Printf("-> [DB] Connecting to MASTER database (%s)...\n", cfg.DSN)

	db, err := sql.Open("sqlite", cfg.DSN)
	if err != nil {
		return nil, nil, err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS books")
	if err == nil {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY, title TEXT, author TEXT)`)
	}
	if err == nil {
		_, err = db.Exec(`INSERT INTO books (title, author) VALUES ('The Go Programming Language', 'Alan A. A. Donovan')`)
	}
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		fmt.Println("-> [DB] Closing MASTER database connection...")
		db.Close()
	}

	return db, cleanup, nil
}

// ProvideReplicaDB provides the Replica database (intended for read operations).
// flora:name=replicaDB
func (c *DBConfig) ProvideReplicaDB(cfg *config.AppConfig) (*sql.DB, func(), error) {
	fmt.Printf("-> [DB] Connecting to REPLICA database (%s)...\n", cfg.DSN)

	db, err := sql.Open("sqlite", cfg.DSN)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		fmt.Println("-> [DB] Closing REPLICA database connection...")
		db.Close()
	}

	return db, cleanup, nil
}
