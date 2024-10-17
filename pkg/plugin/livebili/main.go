package main

import (
	"github.com/futabanaobot/fuu.git/pkg/plugin"
	"github.com/futabanaobot/fuu.git/pkg/plugin/livebili/livebili"
)

func NewPlugin() plugin.Plugin {
	return livebili.NewPlugin()
}
