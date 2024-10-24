package app

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/kohmebot/kohme/pkg/chain"
	"github.com/kohmebot/kohme/pkg/command"
	"github.com/kohmebot/kohme/pkg/gopool"
	"github.com/kohmebot/kohme/pkg/version"
	"github.com/kohmebot/plugin"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
	"strings"
	"time"
)

var v = version.NewVersion(0, 0, 40)

type CoreConf struct {
	HelpTop  string `yaml:"help_top" mapstructure:"help_top"`
	HelpTail string `yaml:"help_tail" mapstructure:"help_tail"`
}

type Core struct {
	app  *App
	conf CoreConf
	db   *gorm.DB
	env  plugin.Env
}

func newCore(a *App) *Core {
	return &Core{
		app: a,
	}
}

func (c *Core) Init(engine *zero.Engine, env plugin.Env) error {
	c.env = env
	err := env.GetConf(&c.conf)
	if err != nil {
		return err
	}
	c.db, err = env.GetDB()
	if err != nil {
		return err
	}
	err = c.db.AutoMigrate(&PluginRecord{})
	if err != nil {
		return err
	}
	err = c.onHelp(engine, env)
	if err != nil {
		return err
	}
	err = c.onPing(engine, env)
	if err != nil {
		return err
	}
	err = c.onPlugin(engine, env)
	if err != nil {
		return err
	}
	err = c.onToggle(engine, env)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) onHelp(engine *zero.Engine, env plugin.Env) error {
	g := env.Groups()
	prefix := c.app.opt.AppConf.Zero.CommandPrefix
	engine.OnCommandGroup([]string{"help", "?", "？", "帮助"}, g.Rule()).Handle(func(ctx *zero.Ctx) {
		var msgChain chain.MessageChain
		msgChain.Split(message.Text(c.conf.HelpTop), message.Text(fmt.Sprintf(`命令前缀 "%s"`, prefix)))
		msgChain.Line()
		for _, name := range c.app.pluginNameSeq {
			pEnv := c.app.envMp[name]
			// 跳过关闭的插件
			if pEnv.disable.Load() {
				continue
			}
			p := c.app.pluginMp[name]
			msgChain.Line(message.Text(fmt.Sprintf("🌟%s (%s)", p.Name(), p.Description())))
			msgChain.Join(message.Text(p.Commands().String()))
		}
		msgChain.Split(message.Text("-----"), message.Text(c.conf.HelpTail))
		gopool.Go(func() {
			ctx.Send(msgChain)
		})
	})
	return nil
}

func (c *Core) onPing(engine *zero.Engine, env plugin.Env) error {
	supers := env.SuperUser()
	engine.OnCommand("ping", supers.Rule()).Handle(func(ctx *zero.Ctx) {
		gopool.Go(func() {
			ctx.Send(message.Text("pong!我还活着"))
		})
	})
	return nil
}

func (c *Core) onPlugin(engine *zero.Engine, env plugin.Env) error {
	supers := env.SuperUser()
	engine.OnCommand("plugin", supers.Rule()).Handle(func(ctx *zero.Ctx) {
		var msgChain chain.MessageChain
		msgChain.Line(message.Text("当前插件列表:"))
		for _, name := range c.app.pluginNameSeq {
			p := c.app.pluginMp[name]
			e := c.app.envMp[name]
			var toggle string
			disable := e.disable.Load()
			if disable {
				toggle = "关闭"
			} else {
				toggle = "开启"
			}
			msgChain.Join(message.Text(fmt.Sprintf("%s v%s (%s)", p.Name(), version.Version(p.Version()).String(), toggle)))
			msgChain.Line()
		}
		gopool.Go(func() {
			ctx.Send(msgChain)
		})
	})
	return nil
}

func (c *Core) onToggle(engine *zero.Engine, env plugin.Env) error {
	supers := env.SuperUser()
	engine.OnCommand("toggle", supers.Rule()).Handle(func(ctx *zero.Ctx) {
		var cmd extension.CommandModel
		var err error
		defer func() {
			if err != nil {
				env.Error(ctx, err)
				return
			}
		}()
		err = ctx.Parse(&cmd)
		if err != nil {
			return
		}
		pluginName := cmd.Args
		pluginName = strings.TrimSpace(pluginName)
		if len(pluginName) <= 0 {
			err = fmt.Errorf("插件名称为空")
			return
		}
		e, ok := c.app.envMp[pluginName]
		if !ok {
			err = fmt.Errorf("插件%s不存在", pluginName)
			return
		}
		var msgChain chain.MessageChain
		if e.disable.CompareAndSwap(true, false) {
			msgChain.SplitEmpty(message.Text(pluginName), message.Text("已开启"))
		} else {
			e.disable.CompareAndSwap(false, true)
			msgChain.SplitEmpty(message.Text(pluginName), message.Text("已关闭"))
		}
		gopool.Go(func() {
			ctx.Send(msgChain)
		})

	})
	return nil
}

func (c *Core) Name() string {
	return "core"
}

func (c *Core) Description() string {
	return "基础插件"
}

func (c *Core) Commands() fmt.Stringer {
	return command.NewCommands(
		command.NewCommand("查看帮助", "help", "?", "帮助"),
	)
}

func (c *Core) Version() uint64 {
	return uint64(v)
}

func (c *Core) OnBoot() {
	var err error
	defer func() {
		if err != nil {
			logrus.Errorf("查询插件校验错误: %s", err.Error())
		}
	}()
	initPluginSet := mapset.NewSet[string]()
	initPluginSet.Append(c.app.pluginNameSeq...)
	var records []PluginRecord
	if err = c.db.Find(&records).Error; err != nil {
		return
	}
	historyPluginMp := make(map[string]PluginRecord, len(records))
	historyPluginSet := mapset.NewSet[string]()
	for _, record := range records {
		historyPluginMp[record.Name] = record
		historyPluginSet.Add(record.Name)
	}
	// 查看是否有新加载的插件
	var newPlugins []plugin.Plugin
	initPluginSet.Difference(historyPluginSet).Each(func(s string) bool {
		newPlugins = append(newPlugins, c.app.pluginMp[s])
		// 返回false才是继续迭代
		return false
	})

	// 查看是否有卸载的插件
	var deletePlugins []string // 卸载的插件只能用string表示
	historyPluginSet.Difference(initPluginSet).Each(func(s string) bool {
		deletePlugins = append(deletePlugins, s)
		return false
	})

	// 查看是否有版本变动的插件
	var updatePlugins []plugin.Plugin
	initPluginSet.Intersect(historyPluginSet).Each(func(s string) bool {
		r := historyPluginMp[s]
		if r.Version != c.app.pluginMp[s].Version() {
			updatePlugins = append(updatePlugins, c.app.pluginMp[s])
		}
		return false
	})
	err = c.db.Transaction(func(tx *gorm.DB) error {
		// 删除记录中已卸载的插件
		if len(deletePlugins) > 0 {
			if err = c.db.Where("name IN ?", deletePlugins).Delete(&PluginRecord{}).Error; err != nil {
				return err
			}
		}
		if len(newPlugins) > 0 {
			// 插入新的插件
			if err = c.db.Create(PluginsToRecord(newPlugins)).Error; err != nil {
				return err
			}
		}
		if len(updatePlugins) > 0 {
			// 更新插件版本
			for _, record := range PluginsToRecord(updatePlugins) {
				if err = c.db.Save(&record).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[%s]\nKohmeBot已启动\n", time.Now().Format("2006-01-02 15:04:05")))
	if len(c.app.pluginNameSeq) > 0 {
		builder.WriteString("已加载插件:\n")
		for idx, s := range c.app.pluginNameSeq {
			p := c.app.pluginMp[s]
			builder.WriteString(fmt.Sprintf("(%d) [%s] v%s\n", idx+1, p.Name(), version.Version(p.Version())))
		}
	}
	if len(newPlugins) > 0 {
		builder.WriteString("新插件:\n")
		for _, p := range newPlugins {
			builder.WriteString(fmt.Sprintf("[%s] v%s\n", p.Name(), version.Version(p.Version())))
		}
	}
	if len(deletePlugins) > 0 {
		builder.WriteString("卸载插件:\n")
		for _, s := range deletePlugins {
			r := historyPluginMp[s]
			builder.WriteString(fmt.Sprintf("[%s] v%s\n", r.Name, version.Version(r.Version)))
		}
	}
	if len(updatePlugins) > 0 {
		builder.WriteString("版本变动:\n")
		for _, p := range updatePlugins {
			hp := historyPluginMp[p.Name()]
			var w string
			if p.Version() > hp.Version {
				w = "版本更新"
			} else {
				w = "版本回退"
			}
			builder.WriteString(fmt.Sprintf("[%s] %s v%s -> v%s\n", p.Name(), w, version.Version(hp.Version), version.Version(p.Version())))
		}
	}
	logrus.Info(builder.String())
	msg := message.Text(builder.String())
	for ctx := range c.env.RangeBot {
		for u := range c.env.SuperUser().RangeUser {
			gopool.Go(func() {
				ctx.SendPrivateMessage(u, msg)
			})
		}
	}

}
