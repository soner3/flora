/*
Copyright © 2026 Soner Astan astansoner@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package errs

import (
	"fmt"
	"hash/fnv"
	"runtime/debug"
	"time"
)

type WeldError struct {
	ID         uint64
	Inner      error
	Message    string
	StackTrace string
	CreatedAt  time.Time
	Misc       map[string]any
}

func Wrap(err error, message string, args ...any) *WeldError {
	formattedMsg := message
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(message, args...)
	}

	if e, ok := err.(*WeldError); ok {
		return &WeldError{
			Inner:      err,
			Message:    formattedMsg,
			StackTrace: e.StackTrace,
			CreatedAt:  time.Now(),
			ID:         e.ID,
			Misc:       nil,
		}
	}

	trace := debug.Stack()

	return &WeldError{
		Inner:      err,
		Message:    formattedMsg,
		StackTrace: string(trace),
		CreatedAt:  time.Now(),
		ID:         GenerateHash(trace),
		Misc:       nil,
	}
}

func (e *WeldError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Inner)
	}
	return e.Message
}

func (e *WeldError) Unwrap() error {
	return e.Inner
}

func GenerateHash(input []byte) uint64 {
	h := fnv.New64a()
	h.Write(input)
	return h.Sum64()
}
