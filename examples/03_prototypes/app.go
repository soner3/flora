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

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/soner3/flora"
)

// =========================================================
// EXAMPLE 1: Simple Prototype (via flora.Component)
// =========================================================

type Session struct {
	// The struct tag changes the lifecycle from singleton to prototype
	flora.Component `flora:"scope=prototype"`
	ID              int
}

func NewSession() *Session {
	id := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000)
	fmt.Printf("   -> [Factory] Created fresh Session [%d]\n", id)
	return &Session{ID: id}
}

// =========================================================
// EXAMPLE 2: Advanced Prototype (via flora.Configuration)
// Useful for request-scoped resources that need cleanup/errors
// =========================================================

// An external struct (e.g. from a database library)
type DBTransaction struct {
	TxID int
}

type DatabaseConfig struct {
	flora.Configuration
}

// flora:scope=prototype
func (c *DatabaseConfig) ProvideTransaction() (*DBTransaction, func(), error) {
	txID := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(999)
	fmt.Printf("   -> [Factory] Opened DB Transaction [%d]\n", txID)

	tx := &DBTransaction{TxID: txID}

	// This cleanup function is returned to the caller of the factory!
	cleanup := func() {
		fmt.Printf("   <- [Cleanup] Closed DB Transaction [%d]\n", txID)
	}

	// We can safely return initialization errors here
	return tx, cleanup, nil
}

// =========================================================
// 3. The Singleton Consumer
// =========================================================

type Server struct {
	flora.Component // Default scope is Singleton (only one server exists)

	// We inject the factory closures instead of the instances!
	sessionFactory func() *Session
	txFactory      func() (*DBTransaction, func(), error)
}

// Flora sees the factory signatures and automatically injects the generated closures.
func NewServer(sFac func() *Session, txFac func() (*DBTransaction, func(), error)) *Server {
	fmt.Println("-> [Init] Server initialized (Singleton)")
	return &Server{
		sessionFactory: sFac,
		txFactory:      txFac,
	}
}

func (s *Server) HandleRequest(user string) {
	fmt.Printf("Server handling request for: %s\n", user)

	// 1. Get a simple prototype instance
	session := s.sessionFactory()

	// 2. Get an advanced prototype instance (with error and cleanup)
	tx, cleanupTx, err := s.txFactory()
	if err != nil {
		fmt.Printf("   [Error] Could not open transaction: %v\n", err)
		return
	}

	// Ensure the prototype resource is safely cleaned up after the request!
	defer cleanupTx()

	fmt.Printf("   Processing %s's data with Session %d and Tx %d\n\n", user, session.ID, tx.TxID)
}
