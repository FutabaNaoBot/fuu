package plugin

import (
	"errors"
	"fmt"
	"github.com/futabanaobot/fuu/internal/util"
	"github.com/futabanaobot/plugin"
	"github.com/sirupsen/logrus"
	"io/fs"
	"path/filepath"
	goplugin "plugin"
)

type Manager struct {
	dir string
}

func NewPluginManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) LoadPlugins() ([]plugin.Plugin, error) {
	if !util.PathExists(m.dir) {
		logrus.Warnf("插件目录%s不存在，跳过加载", m.dir)
		return nil, nil
	}

	var plugins []plugin.Plugin
	err := filepath.Walk(m.dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".so" {
			// 加载单个插件
			plug, err := m.load(path)
			if err != nil {
				return err
			}
			plugins = append(plugins, plug)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("加载插件失败: %w", err)
	}
	return plugins, nil
}

func (m *Manager) LoadPluginsWithFilter(fs ...Filter) ([]plugin.Plugin, error) {
	var res []plugin.Plugin
	plugins, err := m.LoadPlugins()
	if err != nil {
		return nil, err
	}
	for _, plugin := range plugins {
		for _, f := range fs {
			if f(plugin) {
				res = append(res, plugin)
			}
		}
	}
	return res, nil
}

func (m *Manager) load(pluginPath string) (plugin.Plugin, error) {
	p, err := goplugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	newer, err := p.Lookup("NewPlugin")
	if err != nil {
		return nil, errors.New("找不到 NewPlugin")
	}

	getter, ok := newer.(plugin.NewPluginFunc)
	if !ok {
		return nil, fmt.Errorf("NewPlugin( func() Plugin )类型错误,不是%T", getter)
	}
	return getter(), nil

}
