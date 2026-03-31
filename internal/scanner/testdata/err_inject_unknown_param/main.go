package errinjectunknownparam

import "github.com/soner3/flora"

type MyService struct {
	flora.Component `flora:"inject(wrongParamName=masterDB)"`
}

func NewMyService(db string) *MyService { return nil }
