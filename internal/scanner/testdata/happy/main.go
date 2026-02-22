package happy

import "github.com/soner3/weld"

type Iface interface{ Do() }

type A struct {
	weld.Component `weld:""`
}

func NewA() *A { return nil }

type B struct {
	weld.Component `weld:"primary,scope=prototype,constructor=BuildB"`
}

func BuildB() *B { return nil }
func (b *B) Do() {}

type C struct {
	weld.Component `weld:"NewC"`
}

func NewC() *C   { return nil }
func (c *C) Do() {}

type Consumer struct{ weld.Component }

func NewConsumer(i Iface) *Consumer { return nil }
