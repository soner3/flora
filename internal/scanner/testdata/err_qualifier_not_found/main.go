package errqualifiernotfound

import "github.com/soner3/flora"

type MyService struct {
	flora.Component `flora:"inject(db=doesNotExist)"`
}

func NewMyService(db string) *MyService { return nil }
