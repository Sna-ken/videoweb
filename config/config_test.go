package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadReadsConfigFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	content := []byte(`
authorization:
  jwt_secret: test-secret
mysql:
  username: test-user
  password: test-password
  host: mysql.test
  port: "13306"
  name: test-db
  max_open_conns: 30
  max_idle_conns: 15
  conn_max_lifetime: 45m
redis:
  host: redis.test
  port: "16379"
  password: redis-password
  db: 2
rpc:
  auth_address: auth.test:20001
  user_address: user.test:20002
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.MySQL.Host != "mysql.test" {
		t.Fatalf("unexpected mysql host: %q", cfg.MySQL.Host)
	}
	if cfg.MySQL.ConnMaxLifetime != 45*time.Minute {
		t.Fatalf("unexpected mysql connection max lifetime: %s", cfg.MySQL.ConnMaxLifetime)
	}
	if cfg.Redis.DB != 2 {
		t.Fatalf("unexpected redis db: %d", cfg.Redis.DB)
	}
	if cfg.RPC.AuthAddress != "auth.test:20001" {
		t.Fatalf("unexpected auth RPC address: %q", cfg.RPC.AuthAddress)
	}
	if cfg.Authorization.JWTSecret != "test-secret" {
		t.Fatalf("unexpected JWT secret: %q", cfg.Authorization.JWTSecret)
	}

	cachedCfg, err := Load(path)
	if err != nil {
		t.Fatalf("load cached config: %v", err)
	}
	if cachedCfg != cfg {
		t.Fatal("Load returned a different config instance")
	}

	loadedCfg, err := Get()
	if err != nil {
		t.Fatalf("get loaded config: %v", err)
	}
	if loadedCfg != cfg {
		t.Fatal("Get returned a different config instance")
	}

	authorizationCfg, err := GetAuthorization()
	if err != nil {
		t.Fatalf("get authorization config: %v", err)
	}
	if authorizationCfg != &cfg.Authorization {
		t.Fatal("GetAuthorization returned a different authorization config instance")
	}
}
