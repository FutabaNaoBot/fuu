package app

import (
	"fmt"
	"github.com/futabanaobot/fuu.git/pkg/plugin"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

type App struct {
	conf    AConf
	Engine  *zero.Engine
	manager *plugin.Manager

	pluginMp map[string]plugin.Plugin
	envMp    map[string]*Env

	pluginConf PluginConf
}

func New(config AConf, opts ...Option) *App {
	a := &App{
		conf:     config,
		Engine:   zero.New(),
		manager:  plugin.NewPluginManager(config.PluginPath),
		pluginMp: make(map[string]plugin.Plugin),
		envMp:    make(map[string]*Env),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *App) Start() error {
	ps, err := a.manager.LoadPlugins()
	if err != nil {
		return err
	}
	a.AddPlugin(ps...)

	for _, p := range a.pluginMp {
		err = p.Init(a.Engine, a.envMp[p.Name()])
		if err != nil {
			return fmt.Errorf("%s 初始化失败: %w", p.Name(), err)
		}
	}

	zero.RunAndBlock(&a.conf.Zero, func() {
		a.PrintPlugins()
	})
	return nil

}

func (a *App) AddPlugin(ps ...plugin.Plugin) {
	for _, p := range ps {
		a.pluginMp[p.Name()] = p
		pg, ok := a.pluginConf.Plugins[p.Name()]
		if !ok {
			pg = make(map[string]any)
		}
		a.envMp[p.Name()] = NewEnv(p, pg)
	}
}

func (a *App) PrintPlugins() {
	for _, p := range a.pluginMp {
		logrus.Infof("插件 %s | 版本 %s | 描述 %s", p.Name(), p.Version(), p.Description())

	}
}
