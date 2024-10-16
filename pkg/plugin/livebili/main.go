package main

import (
	"botgo/pkg/plugin"
	"botgo/pkg/plugin/livebili/livebili"
)

func NewPlugin() plugin.Plugin {
	return livebili.NewPlugin()
}
