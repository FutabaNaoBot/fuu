package app

import (
	"encoding/json"
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"gopkg.in/yaml.v3"
	"os"
)

type AConf struct {
	PluginPath string `json:"plugin_path"`

	zero.Config
}

func (c *AConf) ParseJsonFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("无法打开 JSON 文件: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(c); err != nil {
		return fmt.Errorf("解析 JSON 文件错误: %w", err)
	}

	return nil
}

type PluginConf struct {
	Plugins map[string]map[string]any `yaml:"plugins"`
}

func (c *PluginConf) ParseYamlFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("无法读取 YAML 文件: %w", err)
	}

	if err := yaml.Unmarshal(file, c); err != nil {
		return fmt.Errorf("解析 YAML 文件错误: %w", err)
	}

	return nil
}
