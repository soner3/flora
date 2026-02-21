package errs

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrap(t *testing.T) {
	testCases := []struct {
		name    string
		err     error
		message string
		args    []any
	}{

		{
			name:    "Test_AlreadyWrappedError",
			err:     Wrap(errors.New("test"), "test"),
			message: "test: %s",
			args:    []any{"test"},
		},
		{
			name:    "Test_NewError",
			err:     errors.New("test"),
			message: "test",
			args:    nil,
		},
	}

	for _, tc := range testCases {
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
