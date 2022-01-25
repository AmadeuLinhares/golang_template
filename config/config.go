package config

import (
	authClient "github.com/tradersclub/TCAuth/client"
	"github.com/tradersclub/TCUtils/cache"
)

// Docs is a struct to use in config
type Docs struct {
	Key string `mapstructure:"key"`
}

// NoSqlDBConn is a struct to use in noSql Database
type NoSqlDBConn struct {
	URL    string `mapstructure:"url"`
	Scheme string `mapstructure:"scheme"`
}

// Server is a struct to use in config
type Server struct {
	Port string `mapstructure:"port"`
}

// DBConn is a struct to use in Database
type DBConn struct {
	URL string `mapstructure:"url"`
}

// Database is a struct to use in config
type Database struct {
	Reader      DBConn      `mapstructure:"reader"`
	Writer      DBConn      `mapstructure:"writer"`
	ReaderNoSql NoSqlDBConn `mapstructure:"reader_no_sql"`
	WriterNoSql NoSqlDBConn `mapstructure:"writer_no_sql"`
}

// Nats is a struct to use in config
type Nats struct {
	URL string `mapstructure:"url"`
}

// TCAuth is a instance to valid Session
type TCAuth struct {
	Addr    string `mapstructure:"addr"`
	Port    string `mapstructure:"port"`
	Timeout string `mapstructure:"timeout"`
}

// Config is a struct to use in var ConfigGlobal
type Config struct {
	ENV      string            `mapstructure:"tc"`
	Docs     Docs              `mapstructure:"docs"`
	Server   Server            `mapstructure:"server"`
	Database Database          `mapstructure:"database"`
	Cache    cache.Options     `mapstructure:"cache"`
	Nats     Nats              `mapstructure:"nats"`
	Auth     authClient.Option `mapstructure:"auth"`
}

// ConfigGlobal is you use in all app
var ConfigGlobal *Config
