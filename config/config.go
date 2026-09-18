package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
	RPC   RPCConfig   `yaml:"rpc"`
}

type MySQLConfig struct {
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	Name            string        `yaml:"name"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type RPCConfig struct {
	AuthAddress string `yaml:"auth_address"`
	UserAddress string `yaml:"user_address"`
}

// Load 读取yml配置文件，应用环境变量覆盖，并验证最终配置
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path is required")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file %q: %w", path, err)
	}
	defer file.Close()

	cfg := Config{
		MySQL: MySQLConfig{
			Host:            "127.0.0.1",
			Port:            "3306",
			MaxOpenConns:    20,
			MaxIdleConns:    10,
			ConnMaxLifetime: 30 * time.Minute,
		},
		Redis: RedisConfig{
			Host: "127.0.0.1",
			Port: "6379",
		},
	}

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config file %q: %w", path, err)
	}

	if err := cfg.applyEnvironment(); err != nil {
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) applyEnvironment() error {
	overrideString("MYSQL_USERNAME", &c.MySQL.Username)
	overrideString("MYSQL_PASSWORD", &c.MySQL.Password)
	overrideString("MYSQL_HOST", &c.MySQL.Host)
	overrideString("MYSQL_PORT", &c.MySQL.Port)
	overrideString("MYSQL_NAME", &c.MySQL.Name)
	overrideString("REDIS_HOST", &c.Redis.Host)
	overrideString("REDIS_PORT", &c.Redis.Port)
	overrideString("REDIS_PASSWORD", &c.Redis.Password)
	overrideString("AUTH_RPC_ADDRESS", &c.RPC.AuthAddress)
	overrideString("USER_RPC_ADDRESS", &c.RPC.UserAddress)

	if err := overrideInt("MYSQL_MAX_OPEN_CONNS", &c.MySQL.MaxOpenConns); err != nil {
		return err
	}
	if err := overrideInt("MYSQL_MAX_IDLE_CONNS", &c.MySQL.MaxIdleConns); err != nil {
		return err
	}
	if err := overrideDuration("MYSQL_CONN_MAX_LIFETIME", &c.MySQL.ConnMaxLifetime); err != nil {
		return err
	}
	if err := overrideInt("REDIS_DB", &c.Redis.DB); err != nil {
		return err
	}

	return nil
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

	return nil
}

func overrideString(name string, target *string) {
	if value, ok := os.LookupEnv(name); ok {
		*target = value
	}
}

func overrideInt(name string, target *int) error {
	value, ok := os.LookupEnv(name)
	if !ok {
		return nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("parse environment variable %s: %w", name, err)
	}
	*target = parsed
	return nil
}

func overrideDuration(name string, target *time.Duration) error {
	value, ok := os.LookupEnv(name)
	if !ok {
		return nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("parse environment variable %s: %w", name, err)
	}
	*target = parsed
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
