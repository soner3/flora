package errwrongtype

import "github.com/soner3/weld"

type Bad struct{ weld.Component }

func NewBad() string { return "" }
