package app

import "botgo/pkg/plugin"

type Option func(app *App)

func WithPlugin(p ...plugin.Plugin) Option {
	return func(app *App) {
		app.AddPlugin(p...)
	}
}

func WithPluginConf(conf PluginConf) Option {
	return func(app *App) {
		app.pluginConf = conf
	}
}
