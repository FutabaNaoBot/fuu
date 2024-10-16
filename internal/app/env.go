package app

import (
	"botgo/pkg/plugin"
	"fmt"
	"github.com/mitchellh/mapstructure"
	zero "github.com/wdvxdr1123/ZeroBot"
	"os"
	"path/filepath"
	"sync/atomic"
)

type Env struct {
	kv map[string]any
	p  plugin.Plugin

	disable atomic.Bool

	ids plugin.Config
}

func NewEnv(p plugin.Plugin, kv map[string]any) *Env {
	e := &Env{
		p:  p,
		kv: kv,
	}
	_ = e.GetConf(&e.ids)
	return e
}

func (e *Env) Rule(r zero.Rule) zero.Rule {
	return func(ctx *zero.Ctx) bool {
		if e.disable.Load() {
			return true
		}
		return r(ctx)
	}
}

func (e *Env) Get(key string) any {
	return e.kv[key]
}

func (e *Env) FilePath() (string, error) {
	path := filepath.Join("data", e.p.Name())
	err := os.MkdirAll(path, os.ModePerm)
	return path, err
}

func (e *Env) RangeBot(yield func(ctx *zero.Ctx) bool) {
	var set []int64
	isDisable := true

	if len(e.ids.Disables) > 0 {
		set = e.ids.Disables
		isDisable = true
	}
	if len(e.ids.Enables) > 0 {
		set = e.ids.Enables
		isDisable = false
	}

	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		var found bool
		for _, target := range set {
			if target == id {
				found = true
				break
			}
		}
		if found && isDisable {
			return true
		} else if !found && !isDisable {
			return true
		}
		return yield(ctx)

	})
}

func (e *Env) GetConf(conf any) error {
	v, ok := e.kv["conf"]
	if !ok {
		return fmt.Errorf("未找到配置")
	}
	vv, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("conf配置类型错误")
	}

	if err := mapstructure.Decode(vv, conf); err != nil {
		return fmt.Errorf("解析配置错误: %v", err)
	}
	return nil
}

func (e *Env) Toggle(b bool) {
	e.disable.Store(!b)
}
