package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Project ProjectConfig `mapstructure:"project"`
	Doris   DBConfig      `mapstructure:"doris"`
	Logger  LogConfig     `mapstructure:"logger"`
	AI      AIConfig      `mapstructure:"ai"`
}

type AIConfig struct {
	APIKey   string `mapstructure:"api_key"`
	BaseURL  string `mapstructure:"base_url"`
	Model    string `mapstructure:"model"`
	MockMode bool   `mapstructure:"mock_mode"`
}

type ProjectConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DBConfig struct {
	Host            string `mapstructure:"host"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Port            int    `mapstructure:"port"`
	HttpPort        int    `mapstructure:"httpport"`
	Database        string `mapstructure:"database"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parse_time"`
	Loc             string `mapstructure:"loc"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

func (l LogConfig) GetLevel() string    { return l.Level }
func (l LogConfig) GetFilename() string { return l.Filename }
func (l LogConfig) GetMaxSize() int     { return l.MaxSize }
func (l LogConfig) GetMaxBackups() int  { return l.MaxBackups }
func (l LogConfig) GetMaxAge() int      { return l.MaxAge }
func (l LogConfig) GetCompress() bool   { return l.Compress }

var GlobalConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		wd, _ := os.Getwd()
		viper.AddConfigPath(filepath.Join(wd, "config"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	GlobalConfig = &cfg
	return GlobalConfig, nil
}
