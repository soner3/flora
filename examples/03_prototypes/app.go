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

// ---------------------------------------------------------
// 1. The Prototype Component
// ---------------------------------------------------------
type Session struct {
	// The struct tag changes the lifecycle from singleton to prototype
	flora.Component `flora:"scope=prototype"`
	ID              int
}

// Every time the factory is called, this constructor runs again
func NewSession() *Session {
	// We generate a random ID to prove it's a fresh instance every time
	id := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000)
	fmt.Printf("-> [Factory] Created a brand new Session with ID: %d\n", id)
	return &Session{ID: id}
}

// ---------------------------------------------------------
// 2. The Singleton Consumer
// ---------------------------------------------------------
type Server struct {
	flora.Component // Default scope is Singleton (only one server exists)

	// IMPORTANT: Because Session is a prototype, we don't inject '*Session'.
	// We inject a factory function that RETURNS a '*Session'!
	sessionFactory func() *Session
}

// Flora sees the factory signature and automatically injects the generated closure.
func NewServer(factory func() *Session) *Server {
	fmt.Println("-> [Init] Server initialized (Singleton)")
	return &Server{sessionFactory: factory}
}

func (s *Server) HandleRequest(user string) {
	// Call the factory to get a fresh session just for this request
	session := s.sessionFactory()
	fmt.Printf("Server handling request for %s using Session %d\n\n", user, session.ID)
}
