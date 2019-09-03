package internal

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// ENV 环境的配置，ID前缀和转发域名
type ENV struct {
	ID     string `toml:"id"`
	Length int    `toml:"length"`
	URL    string `toml:"url"`
}

// ENVS 数组
type ENVS []ENV

// Strategy 策略
type Strategy string

var (
	// StrategyLength 策略 ID长度
	StrategyLength Strategy = "length"
	// StrategyPrefix 策略 ID前缀
	StrategyPrefix Strategy = "prefix"
)

// Config 配置文件
type Config struct {
	URI      string   `toml:"uri"`
	Method   string   `toml:"method"`
	Port     string   `toml:"port"`
	Strategy Strategy `toml:"strategy"`
	ENVS     ENVS     `toml:"envs"`
}

var (
	config = Config{}
)

// Parse 解析配置文件
func Parse(filename string) error {
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// FindEnvs 将IM消息From_Account转化为ID，获取转发的域名
func (config Config) FindEnvs(account string) (ENVS, error) {
	var envs ENVS
	for _, env := range config.ENVS {
		switch config.Strategy {
		case StrategyPrefix:
			if strings.HasPrefix(account, env.ID) {
				envs = append(envs, env)
			}
		case StrategyLength:
			if len(account) == env.Length {
				envs = append(envs, env)
			}
		}
	}
	if len(envs) > 0 {
		return envs, nil
	}
	return nil, fmt.Errorf("not found env config")
}

// GetConfig 获取配置文件
func GetConfig() Config {
	return config
}
