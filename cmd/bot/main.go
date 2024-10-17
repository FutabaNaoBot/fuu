package main

import (
	"github.com/kohmebot/kohme/internal/app"
	"github.com/kohmebot/plugin"
)

var defaultPlugins []plugin.Plugin

func main() {
	conf := app.AConf{}

	err := conf.ParseJsonFile("./conf/config.json")
	if err != nil {
		panic(err)
	}

	pluginConf := app.PluginConf{
		Plugins: map[string]map[string]any{
			"BiliBili-Live": {
				"conf": map[string]any{
					"check_duration": 100,
				},
			},
		}}

	err = pluginConf.ParseYamlFile("./conf/plugins.yaml")
	if err != nil {
		panic(err)
	}

	a := app.New(
		app.WithAppConf(conf),
		app.WithPluginConf(pluginConf),
		app.WithPlugin(defaultPlugins...),
	)
	panic(a.Start())
}

func init() {
	defaultPlugins = plugins()
}
