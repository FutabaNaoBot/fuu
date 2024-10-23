package app

import (
	"fmt"
	"github.com/kohmebot/kohme/internal/db"
	"github.com/kohmebot/plugin"
	"github.com/kohmebot/plugin/pkg/chain"
	"github.com/kohmebot/plugin/pkg/gopool"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sync/atomic"
)

type Env struct {
	customConf   CustomPluginConf
	p            plugin.Plugin
	otherPlugins map[string]plugin.Plugin
	disable      atomic.Bool
	superUser    Users
	group        *GroupsWithEnv
}

func NewEnv(p plugin.Plugin, customConf CustomPluginConf, otherPlugins map[string]plugin.Plugin) *Env {
	e := &Env{
		p:            p,
		customConf:   customConf,
		otherPlugins: otherPlugins,
	}
	e.disable.Store(customConf.Disable)
	e.superUser = customConf.SuperUsers
	e.group = NewGroupsWithEnv(customConf.Groups, e)
	return e
}

func (e *Env) Groups() plugin.Groups {
	return e.group
}

func (e *Env) SuperUser() plugin.Users {
	return e.superUser
}

func (e *Env) Error(ctx *zero.Ctx, err error) {
	if err == nil {
		return
	}
	logrus.Errorf(fmt.Sprintf("[%s] %v", e.p.Name(), err))
	var msgChain chain.MessageChain

	sendToSuperUsers := func() {
		for user := range e.superUser.RangeUser {
			ctx.SendPrivateMessage(user, msgChain)
		}
	}

	send := sendToSuperUsers // 默认情况下发送给超级用户

	if ctx.Event != nil {
		if zero.OnlyGroup(ctx) {
			// 在群聊中需要reply
			msgId := ctx.Event.MessageID
			msgChain.Join(message.Reply(msgId))
		}

		if ctx.Event.GroupID > 0 {
			gid := ctx.Event.GroupID
			send = func() {
				ctx.SendGroupMessage(gid, msgChain)
			}
		} else if ctx.Event.UserID > 0 {
			uid := ctx.Event.UserID
			send = func() {
				ctx.SendPrivateMessage(uid, msgChain)
			}
		}
	}

	msgChain.Split(
		message.Text(fmt.Sprintf("Oops！%s发生错误了！", e.p.Name())),
		message.Text(err.Error()),
	)

	gopool.Go(send)
}

func (e *Env) Get(key string) any {
	return e.customConf.Other[key]
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
	if err := mapstructure.Decode(e.customConf.Conf, conf); err != nil {
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

func (e *Env) GetPlugin(name string) (p plugin.Plugin, ok bool) {
	p, ok = e.otherPlugins[name]
	return
}

func (e *Env) Toggle(b bool) {
	e.disable.Store(!b)
}
