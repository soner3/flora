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
	"errors"
	"fmt"
	"testing"
)

func TestWrap(t *testing.T) {
	testcases := []struct {
		name    string
		err     error
		message string
		args    []any
	}{
		{
			name:    "TestWrapWithWrappedError",
			err:     Wrap(errors.New("test"), "test"),
			message: "test: %s",
			args:    []any{"test"},
		},
		{
			name:    "TestWrapNewError",
			err:     errors.New("test"),
			message: "test",
			args:    nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out := Wrap(tc.err, tc.message, tc.args...)
			if fmt.Sprintf(tc.message, tc.args...) != out.Message {
				t.Errorf("expected: %s, got: %s", tc.message, out.Message)
			}
			if !errors.Is(out, tc.err) {
				t.Errorf("%v is not %v", out, tc.err)
			}
		})
	}

}

func TestError(t *testing.T) {
	const msg = "test"
	inner := Wrap(nil, msg)
	wErr := Wrap(inner, msg)
	if wErr.Error() != fmt.Sprintf("%s: %s", msg, inner.Message) {
		t.Errorf("expected: %s, got: %s", fmt.Sprintf("%s: %s", msg, inner.Message), wErr.Error())
	}
}
