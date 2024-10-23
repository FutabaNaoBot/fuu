package main

import (
	"github.com/kohmebot/chatcount/chatcount"

	"github.com/kohmebot/manager/manager"
	"github.com/kohmebot/plugin"
)

// plugins 加载插件
func plugins() []plugin.Plugin {
	return []plugin.Plugin{
		chatcount.NewPlugin(),
		manager.NewPlugin(),
	}
}
