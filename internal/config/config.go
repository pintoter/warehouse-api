package config

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

const (
	configPath = "./configs"
	configName = "main"
	envFile    = "./.env"
)

type HTTP struct {
	Host            string
	Port            string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

func (h *HTTP) GetAddr() string {
	return fmt.Sprintf("%s:%s", h.Host, h.Port)
}

func (h *HTTP) GetReadTimeout() time.Duration {
	return h.ReadTimeout
}

func (h *HTTP) GetWriteTimeout() time.Duration {
	return h.WriteTimeout
}

func (h *HTTP) GetShutdownTimeout() time.Duration {
	return h.ShutdownTimeout
}

type DB struct {
	User            string
	Password        string
	Host            string
	Port            string
	Name            string
	Sslmode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

func (db *DB) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", db.User, db.Password, db.Host, db.Port, db.Name, db.Sslmode)
}

func (db *DB) GetMaxOpenConns() int {
	return db.MaxOpenConns
}

func (db *DB) GetMaxIdleConns() int {
	return db.MaxIdleConns
}

func (db *DB) GetConnMaxIdleTime() time.Duration {
	return db.ConnMaxIdleTime
}

func (db *DB) GetConnMaxLifetime() time.Duration {
	return db.ConnMaxLifetime
}

type Project struct {
	Name  string
	Level string
	Mode  string
}

func (p *Project) GetName() string {
	return p.Name
}

func (p *Project) GetLevel() string {
	return p.Level
}

func (p *Project) GetMode() string {
	return p.Mode
}

type Config struct {
	HTTP
	DB
	Project
}

var config = new(Config)
var once sync.Once

func init() {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("loading env file")
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("reading config err")
	}
}

func New() *Config {
	once.Do(func() {
		var err error

		err = viper.Unmarshal(config)
		if err != nil {
			log.Fatal("reading config")
		}

		err = envconfig.Process("db", &config.DB)
		if err != nil {
			log.Fatal("error: get env for db")
		}
	})
	return config
}
