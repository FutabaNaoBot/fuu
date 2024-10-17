package main

import "github.com/futabanaobot/plugin"

var defaultPlugins []plugin.Plugin

// plugins 加载插件
func plugins() []plugin.Plugin {
	return []plugin.Plugin{}
}

func init() {
	defaultPlugins = plugins()
}
