package app

import (
	"fmt"
	"github.com/kohmebot/plugin"
	"github.com/kohmebot/plugin/pkg/chain"
	"github.com/kohmebot/plugin/pkg/command"
	"github.com/kohmebot/plugin/pkg/gopool"
	"github.com/kohmebot/plugin/pkg/version"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var v = version.NewVersion(0, 0, 1)

type CoreConf struct {
	HelpTop  string `yaml:"help_top" mapstructure:"help_top"`
	HelpTail string `yaml:"help_tail" mapstructure:"help_tail"`
}

type Core struct {
	app  *App
	conf CoreConf
}

func newCore(a *App) *Core {
	return &Core{
		app: a,
	}
}

func (c *Core) Init(engine *zero.Engine, env plugin.Env) error {

	err := env.GetConf(&c.conf)
	if err != nil {
		return err
	}
	err = c.onHelp(engine, env)
	if err != nil {
		return err
	}
	err = c.onManager(engine, env)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) onHelp(engine *zero.Engine, env plugin.Env) error {
	g := env.Groups()
	prefix := c.app.opt.AppConf.Zero.CommandPrefix
	engine.OnCommandGroup([]string{"help", "?", "？", "帮助"}, g.Rule(func(ctx *zero.Ctx) bool {
		var msgChain chain.MessageChain
		msgChain.Split(message.Text(c.conf.HelpTop), message.Text(fmt.Sprintf(`命令前缀 "%s"`, prefix)))
		msgChain.Line()
		for _, p := range c.app.pluginMp {
			pEnv := c.app.envMp[p.Name()]
			// 跳过关闭的插件
			if pEnv.disable.Load() {
				continue
			}

			msgChain.Line(message.Text(fmt.Sprintf("🌟%s (%s)", p.Name(), p.Description())))
			msgChain.Join(message.Text(p.Commands().String()))
		}
		msgChain.Split(message.Text("-----"), message.Text(c.conf.HelpTail))
		gopool.Go(func() {
			ctx.Send(msgChain)
		})

		return true
	}))
	return nil
}

func (c *Core) onManager(engine *zero.Engine, env plugin.Env) error {
	supers := c.app.opt.AppConf.Zero.SuperUsers
	engine.OnCommand("ping", func(ctx *zero.Ctx) bool {
		for _, super := range supers {
			if super == ctx.Event.Sender.ID {
				gopool.Go(func() {
					ctx.Send(message.Text("pong!我还活着"))
				})
				break
			}
		}
		return true
	})
	return nil
}

func (c *Core) Name() string {
	return "core"
}

func (c *Core) Description() string {
	return "核心插件"
}

func (c *Core) Commands() command.Commands {
	return command.NewCommands(
		command.NewCommand("查看帮助", "help", "?", "帮助"),
	)
}

func (c *Core) Version() version.Version {
	return v
}
