package app

import "github.com/kohmebot/plugin"

type Option func(opt *option)

type option struct {
	PluginConf     PluginConf
	AppConf        AConf
	DefaultPlugins []plugin.Plugin
}

func WithPlugin(p ...plugin.Plugin) Option {
	return func(opt *option) {
		opt.DefaultPlugins = append(opt.DefaultPlugins, p...)
	}
}

func WithPluginConf(conf PluginConf) Option {
	return func(opt *option) {
		opt.PluginConf = conf
	}
}

func WithAppConf(conf AConf) Option {
	return func(opt *option) {
		opt.AppConf = conf
	}
}

func defaultOption() option {
	return option{
		PluginConf:     PluginConf{},
		AppConf:        AConf{},
		DefaultPlugins: nil,
	}
}
