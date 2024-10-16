package livebili

import (
	"botgo/pkg/plugin"
	"botgo/pkg/version"
	zero "github.com/wdvxdr1123/ZeroBot"
)

type biliPlugin struct {
	e         *zero.Engine
	env       plugin.Env
	liveState map[int64]bool
	conf      Config
}

func NewPlugin() plugin.Plugin {
	return &biliPlugin{liveState: make(map[int64]bool)}
}

func (b *biliPlugin) Init(engine *zero.Engine, env plugin.Env) error {
	b.e = engine
	b.env = env

	return b.init()
}

func (b *biliPlugin) Name() string {
	return "BiliBili-Live"
}

func (b *biliPlugin) Description() string {
	return "None"
}

func (b *biliPlugin) Help() string {
	return "None"
}

func (b *biliPlugin) Version() version.Version {
	return version.NewVersion(0, 0, 1)
}
