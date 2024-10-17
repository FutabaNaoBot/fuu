package app

import (
	"fmt"
	"github.com/kohmebot/plugin"
	"github.com/kohmebot/plugin/pkg/command"
	"github.com/kohmebot/plugin/pkg/version"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strings"
)

var v = version.NewVersion(0, 0, 1)

type Core struct {
	app *App
}

func newCore(a *App) *Core {
	return &Core{
		a,
	}
}

func (c *Core) Init(engine *zero.Engine, env plugin.Env) error {
	// FIXME 功能待完善。暂时关闭
	g := env.Groups()
	return nil
	engine.OnCommandGroup([]string{"help", "?", "？"}, g.Rule(func(ctx *zero.Ctx) bool {
		var builder strings.Builder
		for _, p := range c.app.pluginMp {
			if p.Name() == c.Name() {
				continue
			}
			builder.WriteString(fmt.Sprintf("%s (%s)\n", p.Name(), p.Description()))
			builder.WriteString(p.Commands().String())
			builder.WriteByte('\n')
		}
		go ctx.Send(message.Text(builder.String()))
		return true
	}))
	return nil
}

func (c *Core) Name() string {
	return "core"
}

func (c *Core) Description() string {
	return "管理装载的插件"
}

func (c *Core) Commands() command.Commands {
	return command.NewCommands(
		command.NewCommand("查看帮助", "help", "?", "？"),
	)
}

func (c *Core) Version() version.Version {
	return v
}
