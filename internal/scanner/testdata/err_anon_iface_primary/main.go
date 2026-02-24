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
package erranonifaceprimary

import "github.com/soner3/flora"

type Impl1 struct {
	flora.Component `flora:"primary"`
}

func NewImpl1() *Impl1 { return nil }
func (i *Impl1) Do()   {}

type Impl2 struct{ flora.Component }

func NewImpl2() *Impl2 { return nil }
func (i *Impl2) Do()   {}

type Bad struct{ flora.Component }

func NewBad(req interface{ Do() }) *Bad { return nil }
