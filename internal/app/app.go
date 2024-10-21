package app

import (
	"fmt"
	fplugin "github.com/kohmebot/kohme/pkg/plugin"
	"github.com/kohmebot/plugin"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

type App struct {
	opt option

	Engine  *zero.Engine
	manager *fplugin.Manager

	pluginMp map[string]plugin.Plugin
	// 插件名称序列，这个表明了加载顺序
	pluginNameSeq []string
	envMp         map[string]*Env
}

func New(opts ...Option) *App {
	defaultOpt := defaultOption()
	for _, opt := range opts {
		opt(&defaultOpt)
	}

	a := &App{
		opt:      defaultOpt,
		Engine:   zero.New(),
		manager:  fplugin.NewPluginManager(defaultOpt.PluginConf.Path),
		pluginMp: make(map[string]plugin.Plugin),
		envMp:    make(map[string]*Env),
	}

	return a
}

func (a *App) Start() error {
	ps, err := a.manager.LoadPlugins()
	if err != nil {
		return err
	}
	a.AddPlugin(newCore(a))
	a.AddPlugin(append(a.opt.DefaultPlugins, ps...)...)

	for _, name := range a.pluginNameSeq {
		p := a.pluginMp[name]
		err = p.Init(a.Engine, a.envMp[p.Name()])
		if err != nil {
			return fmt.Errorf("%s 初始化失败: %w", p.Name(), err)
		}
	}
	a.PrintPlugins()
	zero.RunAndBlock(&a.opt.AppConf.Zero, func() {
		for _, name := range a.pluginNameSeq {
			a.pluginMp[name].OnBoot()
		}
	})
	return nil

}

func (a *App) AddPlugin(ps ...plugin.Plugin) {
	for _, p := range ps {
		if _, ok := a.pluginMp[p.Name()]; ok {
			panic(fmt.Sprintf("存在重复名称的插件: %s", p.Name()))
		}
		a.pluginMp[p.Name()] = p
		a.pluginNameSeq = append(a.pluginNameSeq, p.Name())
		pg, ok := a.opt.PluginConf.Plugins[p.Name()]
		if !ok || pg == nil {
			pg = make(map[string]any)
		}
		_, ok = pg["groups"]
		if !ok {
			pg["groups"] = a.opt.PluginConf.Groups
		}
		pg["super_users"] = a.opt.AppConf.Zero.SuperUsers
		a.envMp[p.Name()] = NewEnv(p, pg)
	}
}

func (a *App) PrintPlugins() {
	for _, name := range a.pluginNameSeq {
		p := a.pluginMp[name]
		logrus.Infof("插件 %s | 版本 %s | 描述 %s", p.Name(), p.Version(), p.Description())
	}
}
