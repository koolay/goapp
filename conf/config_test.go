package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		cfgFile string
	}
	const configFile = "testdata/app.toml"

	appCfg, err := LoadConfig(configFile)
	assert.Nil(t, err)
	assert.NotNil(t, appCfg)
	assert.Equal(t, appCfg.App.Name, "goapp")
	assert.Equal(t, appCfg.DB.Password, "goapp#123")

	// 测试默认配置值
	deafultCfg := defaultConfiguration()
	// string类型
	assert.Equal(t, appCfg.Log.Level, deafultCfg.Log.Level)
	// int类型
	assert.Equal(t, appCfg.DB.Port, deafultCfg.DB.Port)
}

func Test_defaultConfiguration(t *testing.T) {
	cfg := defaultConfiguration()
	assert.Equal(t, cfg.App.Name, "goapp")
}
