package errcollisionnoprimary

import "github.com/soner3/weld"

type Greeter interface {
	Greet()
}
type GreeterA struct {
	weld.Component
}

func NewGreeterA() *GreeterA { return nil }
func (g *GreeterA) Greet()   {}

type GreeterB struct {
	weld.Component
}

func NewGreeterB() *GreeterB { return nil }
func (g *GreeterB) Greet()   {}

type Consumer struct {
	weld.Component
}

func NewConsumer(g Greeter) *Consumer { return nil }
