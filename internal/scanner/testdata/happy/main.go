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
package happy

import "github.com/soner3/flora"

type Iface interface{ Do() }

type A struct {
	flora.Component `flora:""`
}

func NewA() *A { return nil }

type B struct {
	flora.Component `flora:"primary,scope=prototype,constructor=BuildB,order=1"`
}

func BuildB() *B { return nil }
func (b *B) Do() {}

type C struct {
	flora.Component `flora:"NewC"`
}

func NewC() *C   { return nil }
func (c *C) Do() {}

type Consumer struct{ flora.Component }

func NewConsumer(i Iface) *Consumer { return nil }

type ProtoConsumer struct{ flora.Component }

func NewProtoConsumer(ifaceFactory func() Iface, structFactory func() *B) *ProtoConsumer {
	return nil
}
