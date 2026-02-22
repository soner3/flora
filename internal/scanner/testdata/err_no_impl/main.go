package errnoimpl

import "github.com/soner3/weld"

type Iface interface{ Do() }

type Consumer struct{ weld.Component }

func NewConsumer(i Iface) *Consumer { return nil }
