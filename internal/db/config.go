package db

type Config struct {
	Path string
	// 最大空闲连接数
	IdleConn int `yaml:"idle_conn"`
	// 最大连接数
	MaxConn int `yaml:"max_conn"`
	// 连接最大存活时间(分钟)
	MaxLifeTime int64 `yaml:"max_life_time"`
}
