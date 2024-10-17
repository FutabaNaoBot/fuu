package app

import (
	"fmt"
	"github.com/kohmebot/kohme/internal/db"
	"github.com/kohmebot/kohme/internal/util"
	"github.com/kohmebot/plugin"
	"github.com/mitchellh/mapstructure"
	zero "github.com/wdvxdr1123/ZeroBot"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sync/atomic"
)

type Env struct {
	kv      map[string]any
	p       plugin.Plugin
	disable atomic.Bool
}

func NewEnv(p plugin.Plugin, kv map[string]any) *Env {
	e := &Env{
		p:  p,
		kv: kv,
	}
	disable, ok := e.Get("disable").(bool)
	if ok {
		e.disable.Store(disable)
	}

	return e
}

func (e *Env) Groups() plugin.Groups {
	vv := e.kv["groups"]
	switch res := vv.(type) {
	case []any:
		i64s, ok := util.AnySliceToInt64(res)
		if !ok {
			break
		}
		return Groups(i64s)
	case Groups:
		return res
	case []int:
		return Groups(util.ToInt64Slice(res))
	case []int32:
		return Groups(util.ToInt64Slice(res))
	case []int64:
		return Groups(res)
	}
	return Groups([]int64{})
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
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
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

func (e *Env) GetDB() (*gorm.DB, error) {
	p, err := e.FilePath()
	if err != nil {
		return nil, err
	}
	return db.Get(filepath.Join(p, fmt.Sprintf("%s.db", e.p.Name())))
}

func (e *Env) Toggle(b bool) {
	e.disable.Store(!b)
}
