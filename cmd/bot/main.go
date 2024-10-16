package main

import (
	"botgo/internal/app"
	"botgo/pkg/plugin/livebili/livebili"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	conf := app.AConf{
		PluginPath: "./plugins",
		Zero: zero.Config{
			NickName:      []string{"bot"},
			CommandPrefix: "/",
			SuperUsers:    []int64{123456},
			Driver: []zero.Driver{
				// 正向 WS
				driver.NewWebSocketClient("ws://127.0.0.1:6700", ""),
				// 反向 WS
				driver.NewWebSocketServer(16, "ws://127.0.0.1:6701", ""),
			},
		},
	}

	err := conf.ParseJsonFile("./config.json")
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

	err = pluginConf.ParseYamlFile("./plugins.yaml")
	if err != nil {
		panic(err)
	}

	a := app.New(conf,
		app.WithPluginConf(pluginConf),
		app.WithPlugin(livebili.NewPlugin()),
	)
	panic(a.Start())
}
