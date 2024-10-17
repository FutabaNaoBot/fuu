package plugin

import (
	"github.com/jhue58/botgo/pkg/command"
	"github.com/jhue58/botgo/pkg/version"
	zero "github.com/wdvxdr1123/ZeroBot"
	"gorm.io/gorm"
)

type NewPluginFunc = func() Plugin

// Plugin 所有插件需实现的接口
type Plugin interface {
	Init(engine *zero.Engine, env Env) error
	Name() string
	Description() string
	Commands() command.Commands
	Version() version.Version
}

type Env interface {
	Get(key string) any
	FilePath() (string, error)
	Rule(r zero.Rule) zero.Rule
	GetConf(conf any) error
	GetDB() (*gorm.DB, error)
	RangeBot(yield func(ctx *zero.Ctx) bool)
}
