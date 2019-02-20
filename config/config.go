package config

import (
	"time"

	"github.com/shibukawa/configdir"
	"github.com/spf13/viper"
)

// Provider defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

var defaultConfig *viper.Viper

// ConfigDir is the directory where config is stored
var ConfigDir *configdir.Config

// Config gets the default config
func Config() Provider {
	return defaultConfig
}

// LoadConfigProvider loads the config
func LoadConfigProvider(appName string) Provider {
	return readViperConfig(appName)
}

func init() {
	configDirs := configdir.New("krontab", "krontab")
	folders := configDirs.QueryFolders(configdir.Global)
	ConfigDir = folders[0]
	defaultConfig = readViperConfig("KRONTAB")
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath(ConfigDir.Path)
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	// global defaults
	v.SetDefault("json_logs", false)
	v.SetDefault("loglevel", "debug")
	v.SetDefault("namespace", "default")

	v.ReadInConfig()

	return v
}
