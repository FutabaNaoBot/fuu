package plugin

type Config struct {
	Enables  []int64 `yaml:"enables"`
	Disables []int64 `yaml:"disables"`
}
