package config

import (
	"path/filepath"
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

// SystemConfigDir is the directory where system config is stored
var SystemConfigDir *configdir.Config

// TemplateDirs are the directorys to check for templates
var TemplateDirs []string

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
	ConfigDir = configDirs.QueryFolders(configdir.Global)[0]
	SystemConfigDir = configDirs.QueryFolders(configdir.System)[0]
	TemplateDirs = append(TemplateDirs, filepath.Join(ConfigDir.Path, "templates"))
	TemplateDirs = append(TemplateDirs, filepath.Join(SystemConfigDir.Path, "templates"))

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

	v.ReadInConfig()

	return v
}
