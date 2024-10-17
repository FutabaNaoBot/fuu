package main

import (
	"github.com/jhue58/botgo/pkg/plugin"
	"github.com/jhue58/botgo/pkg/plugin/livebili/livebili"
)

func NewPlugin() plugin.Plugin {
	return livebili.NewPlugin()
}
