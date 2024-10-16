package plugin

import (
	"botgo/pkg/version"
	zero "github.com/wdvxdr1123/ZeroBot"
)

type NewPluginFunc = func() Plugin

// Plugin 所有插件需实现的接口
type Plugin interface {
	Init(engine *zero.Engine, env Env) error
	Name() string
	Description() string
	Help() string
	Version() version.Version
}

type Env interface {
	Get(key string) any
	FilePath() (string, error)
	Rule(r zero.Rule) zero.Rule
	GetConf(conf any) error
	RangeBot(yield func(ctx *zero.Ctx) bool)
}
