package happy

import "github.com/soner3/weld"

type Greeter interface {
	Greet() string
}

type SimpleLogger struct {
	weld.Component
}

func NewSimpleLogger() SimpleLogger {
	return SimpleLogger{}
}

type GermanGreeter struct {
	weld.Component `weld:"constructor=BuildGermanGreeter,"`
}

func BuildGermanGreeter() *GermanGreeter {
	return &GermanGreeter{}
}

func (g *GermanGreeter) Greet() string {
	return "Hallo"
}

type App struct {
	weld.Component
}

func NewApp(g Greeter, l SimpleLogger) *App {
	return &App{}
}

type JustANormalStruct struct {
	SomeConfig string
	Value      int
}

type UntaggedComponent struct {
	weld.Component
}

func NewUntaggedComponent() *UntaggedComponent {
	return nil
}

var _ = "Trigger obj == nil check"

func init() {
	// Dummy init
}
