package config

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL         MySQLConfig         `mapstructure:"mysql" yaml:"mysql"`
	Redis         RedisConfig         `mapstructure:"redis" yaml:"redis"`
	RPC           RPCConfig           `mapstructure:"rpc" yaml:"rpc"`
	Authorization AuthorizationConfig `mapstructure:"authorization" yaml:"authorization"`
}

type MySQLConfig struct {
	Username        string        `mapstructure:"username" yaml:"username"`
	Password        string        `mapstructure:"password" yaml:"password"`
	Host            string        `mapstructure:"host" yaml:"host"`
	Port            string        `mapstructure:"port" yaml:"port"`
	Name            string        `mapstructure:"name" yaml:"name"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"`
}

type RPCConfig struct {
	AuthAddress string `mapstructure:"auth_address" yaml:"auth_address"`
	UserAddress string `mapstructure:"user_address" yaml:"user_address"`
}

type AuthorizationConfig struct {
	JWTSecret string `mapstructure:"jwt_secret" yaml:"jwt_secret"`
}

var (
	loadOnce     sync.Once
	configMu     sync.RWMutex
	loadedConfig *Config
	loadErr      error
)

// Load 在当前进程中只读取和解析一次指定的配置文件。
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path is required")
	}

	loadOnce.Do(func() {
		cfg, err := loadFile(path)
		if err != nil {
			loadErr = err
			return
		}

		configMu.Lock()
		loadedConfig = cfg
		configMu.Unlock()
	})

	if loadErr != nil {
		return nil, loadErr
	}

	return Get()
}

// Get 返回 Load 已初始化的配置实例，不会重新读取配置文件。
func Get() (*Config, error) {
	configMu.RLock()
	defer configMu.RUnlock()

	if loadedConfig == nil {
		return nil, errors.New("config is not loaded")
	}

	return loadedConfig, nil
}

// GetAuthorization 返回已初始化配置中的 Authorization 实例。
func GetAuthorization() (*AuthorizationConfig, error) {
	cfg, err := Get()
	if err != nil {
		return nil, err
	}

	return &cfg.Authorization, nil
}

func loadFile(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file %q: %w", path, err)
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return nil, fmt.Errorf("decode config file %q: %w", path, err)
	}
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.MySQL.Username == "" {
		return errors.New("mysql.username is required")
	}
	if c.MySQL.Host == "" {
		return errors.New("mysql.host is required")
	}
	if err := validatePort("mysql.port", c.MySQL.Port); err != nil {
		return err
	}
	if c.MySQL.Name == "" {
		return errors.New("mysql.name is required")
	}
	if c.MySQL.MaxOpenConns <= 0 {
		return errors.New("mysql.max_open_conns must be greater than 0")
	}
	if c.MySQL.MaxIdleConns < 0 {
		return errors.New("mysql.max_idle_conns must not be negative")
	}
	if c.MySQL.MaxIdleConns > c.MySQL.MaxOpenConns {
		return errors.New("mysql.max_idle_conns must not exceed mysql.max_open_conns")
	}
	if c.MySQL.ConnMaxLifetime <= 0 {
		return errors.New("mysql.conn_max_lifetime must be greater than 0")
	}

	if c.Redis.Host == "" {
		return errors.New("redis.host is required")
	}
	if err := validatePort("redis.port", c.Redis.Port); err != nil {
		return err
	}
	if c.Redis.DB < 0 {
		return errors.New("redis.db must not be negative")
	}
	if c.RPC.AuthAddress == "" {
		return errors.New("rpc.auth_address is required")
	}
	if c.RPC.UserAddress == "" {
		return errors.New("rpc.user_address is required")
	}
	if c.Authorization.JWTSecret == "" {
		return errors.New("authorization.jwt_secret is required")
	}

	return nil
}

func validatePort(name, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%s must be a number: %w", name, err)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("%s must be between 1 and 65535", name)
	}
	return nil
}
