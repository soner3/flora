package errnotfunc

import "github.com/soner3/weld"

type BadComponent struct {
	weld.Component
}

var NewBadComponent = "not a function"
