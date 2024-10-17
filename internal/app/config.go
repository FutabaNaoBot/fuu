package app

import (
	"encoding/json"
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"gopkg.in/yaml.v3"
	"os"
)

type WsConf struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

type AConf struct {
	Zero      zero.Config `json:"zero"`
	Ws        WsConf      `json:"ws"`
	ReverseWs WsConf      `json:"rws"`
}

func (c *AConf) InitDriver() {
	var ds []zero.Driver
	if c.Ws.Url != "" {
		// 正向Ws
		ds = append(ds, driver.NewWebSocketClient(c.Ws.Url, c.Ws.Token))
	}
	if c.ReverseWs.Url != "" {
		// 反向Ws
		ds = append(ds, driver.NewWebSocketServer(16, c.ReverseWs.Url, c.ReverseWs.Token))
	}
	c.Zero.Driver = ds
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

	c.InitDriver()

	return nil
}

type PluginConf struct {
	Path    string                    `yaml:"path"`
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
