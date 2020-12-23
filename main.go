package main

import (
	"fmt"
	"io"
	"os"

	"github.com/caarlos0/env"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Core config structure to unmarshal
type Core struct {
	Delta int `env:"DELTA"`
}

// Log config structure to unmarshal
type Log struct {
	Level uint     `env:"LOG_LEVEL"`
	Path  []string `env:"LOG_PATH"`
	Time  bool     `env:"LOG_TIME"`
}

// Database config structure to unmarshal
type Database struct {
	Host string `env:"DATABASE_HOST"`
	Port string `env:"DATABASE_PORT"`
	Name string `env:"DATABASE_NAME"`
	User string `env:"DATABASE_USER"`
	Pass string `env:"DATABASE_PASS"`
	Zone string `env:"DATABASE_ZONE"`
	Ssl  string `env:"DATABASE_SSL"`
}

// Imap config structure to unmarshal
type Imap struct {
	Host string `env:"IMAP_HOST"`
	Port string `env:"IMAP_PORT"`
	User string `env:"IMAP_USER"`
	Pass string `env:"IMAP_PASS"`
	Msg  int `env:"IMAP_MSG"`
}

// Emailbot config structure to unmarshal
type Emailbot struct {
	Proc int `env:"EMAILBOT_PROC"`
}

// Conf configuration object
type Conf struct {
	viper        *viper.Viper
	confCore     Core
	confLog      Log
	confDatabase Database
	confImap     Imap
	confEmailbot Emailbot
}

// functions to call in constructor

func loadCore(cfg *Conf) error {
	// unmarshal Conf to struct
	if err := cfg.viper.Unmarshal(&cfg.confCore); err != nil {
		return err
	}
	// get log configuration (if set) from env
	_ = env.Parse(&cfg.confCore)

	return nil
}

func loadLog(cfg *Conf) error {
	// unmarshal Conf to struct
	if err := cfg.viper.UnmarshalKey("log", &cfg.confLog); err != nil {
		return err
	}
	// get log configuration (if set) from env
	_ = env.Parse(&cfg.confLog)

	return nil
}

func loadDatabase(cfg *Conf) error {
	// unmarshal Conf to struct
	if err := cfg.viper.UnmarshalKey("database", &cfg.confDatabase); err != nil {
		return err
	}
	// get log configuration (if set) from env
	_ = env.Parse(&cfg.confDatabase)

	return nil
}

func loadImap(cfg *Conf) error {
	// unmarshal Conf to struct
	if err := cfg.viper.UnmarshalKey("imap", &cfg.confImap); err != nil {
		return err
	}
	// get log configuration (if set) from env
	_ = env.Parse(&cfg.confImap)

	return nil
}

func loadEmailbot(cfg *Conf) error {
	// unmarshal Conf to struct
	if err := cfg.viper.UnmarshalKey("emailbot", &cfg.confEmailbot); err != nil {
		return err
	}
	// get log configuration (if set) from env
	_ = env.Parse(&cfg.confEmailbot)

	return nil
}

// NewConfiguration constructor for Conf object
func NewConfiguration() (*Conf,error) {
	var cfg = &Conf{}

	cfg.viper = viper.New()
	// config file name (without extension)
	cfg.viper.SetConfigName("config")
	// config extension
	cfg.viper.SetConfigType("yaml")
	// config file paths
	cfg.viper.AddConfigPath("/etc/app/")
	// when it's called inside main package
	cfg.viper.AddConfigPath(".")
	// read config
	err := cfg.viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if err := loadCore(cfg); err != nil {
		return nil, err
	}
	if err := loadLog(cfg); err != nil {
		return nil, err
	}
	if err := loadDatabase(cfg); err != nil {
		return nil, err
	}
	if err := loadImap(cfg); err != nil {
		return nil, err
	}
	if err := loadEmailbot(cfg); err != nil {
		return nil, err
	}

	// update config file dynamically
	cfg.viper.WatchConfig()
	cfg.viper.OnConfigChange(func(e fsnotify.Event) {
		if e.Op == fsnotify.Write {
			fmt.Println("Config file changed, reload", e.Name)

			if err := loadCore(cfg); err != nil {
				panic(err)
			}
			if err := loadLog(cfg); err != nil {
				panic(err)
			}
			if err := loadDatabase(cfg); err != nil {
				panic(err)
			}
			if err := loadImap(cfg); err != nil {
				panic(err)
			}
			if err := loadEmailbot(cfg); err != nil {
				panic(err)
			}
		}
	})
	return cfg, nil
}

// GetDelta delta time, application cycle wait
func (cfg *Conf) GetDelta() int {
	return cfg.confCore.Delta
}

// GetLogWriters paths of logs
func (cfg *Conf) GetLogWriters() ([]io.Writer,error) {
	var writers []io.Writer
	var writerErr error
	writers = append(writers, os.Stdout)
	for _, path := range cfg.confLog.Path {
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		writerErr = err
		writers = append(writers, writer)
	}
	return writers, writerErr
}

// GetLogLVL the lvl of the log
func (cfg *Conf) GetLogLVL() uint  {
	return cfg.confLog.Level
}

// GetLogTS enable or disable TimeStamp
func (cfg *Conf) GetLogTS() bool  {
	return cfg.confLog.Time
}

// GetDBConnStr conn string of Postgres DB
func (cfg *Conf) GetDBConnStr() string  {
	return fmt.Sprintf("host=%s  port=%s  dbname=%s  user=%s  password=%s  sslmode=%s  TimeZone=%s",
		cfg.confDatabase.Host,
		cfg.confDatabase.Port,
		cfg.confDatabase.Name,
		cfg.confDatabase.User,
		cfg.confDatabase.Pass,
		cfg.confDatabase.Ssl,
		cfg.confDatabase.Zone,
	)
}

// GetImapServer return host:port
func (cfg *Conf) GetImapServer() string {
	return fmt.Sprintf("%s:%s", cfg.confImap.Host,cfg.confImap.Port)
}

// GetImapUser return imap user
func (cfg *Conf) GetImapUser() string {
	return cfg.confImap.User
}

// GetImapPass return imap pass
func (cfg *Conf) GetImapPass() string {
	return cfg.confImap.Pass
}

// GetImapMsg messages at a time (per process)
func (cfg *Conf) GetImapMsg() int {
	return cfg.confImap.Msg
}

// GetEmailbotProc mailbox process at a time
func (cfg *Conf) GetEmailbotProc() int {
	return cfg.confEmailbot.Proc
}


func main(){
	conf, err := NewConfiguration()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", conf.GetDelta())
	writers, err := conf.GetLogWriters()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", writers)
	fmt.Printf("%#v\n", conf.GetLogLVL())
	fmt.Printf("%#v\n", conf.GetLogTS())
	fmt.Printf("%#v\n", conf.GetDBConnStr())
	fmt.Printf("%#v\n", conf.GetImapServer())
	fmt.Printf("%#v\n", conf.GetImapUser())
	fmt.Printf("%#v\n", conf.GetImapPass())
	fmt.Printf("%#v\n", conf.GetEmailbotProc())
}
