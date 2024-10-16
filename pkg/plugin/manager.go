package plugin

import (
	"errors"
	"fmt"
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

func (m *Manager) LoadPlugins() ([]Plugin, error) {
	var plugins []Plugin
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

func (m *Manager) LoadPluginsWithFilter(fs ...Filter) ([]Plugin, error) {
	var res []Plugin
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

func (m *Manager) load(pluginPath string) (Plugin, error) {
	p, err := goplugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	newer, err := p.Lookup("NewPlugin")
	if err != nil {
		return nil, errors.New("找不到 NewPlugin")
	}

	getter, ok := newer.(NewPluginFunc)
	if !ok {
		return nil, fmt.Errorf("NewPlugin( func() Plugin )类型错误,不是%T", getter)
	}
	return getter(), nil

}
