package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"

	"github.com/spf13/viper"
)

var Instance Config

type Config struct {
	Host string `yaml:"Host"`

	// spreadsheet
	SpreadSheetId string `yaml:"SpreadSheetId"`
	SheetName     string `yaml:"SheetName"`
	MaxFetchRange string `yaml:"MaxFetchRange"`
	StatusColumn  string `yaml:"StatusColumn"`
	RemarkColumn  string `yaml:"RemarkColumn"`

	// queue
	DelayWhenNoJobInSeconds   int `yaml:"DelayWhenNoJobInSeconds"`
	DelayWhenErrorInSeconds   int `yaml:"DelayWhenErrorInSeconds"`
	DelayWhenJobDoneInSeconds int `yaml:"DelayWhenJobDoneInSeconds"`

	// google api cred
	CredentialPath string `yaml:"CredentialPath"`
	TokenPath      string `yaml:"TokenPath"`

	// db
	ShowSql      bool        `yaml:"ShowSql"`
	MySqlUrl     string      `yaml:"MySqlUrl"`
	MySqlMaxIdle int         `yaml:"MySqlMaxIdle"`
	MySqlMaxOpen int         `yaml:"MySqlMaxOpen"`
	Redis        RedisConfig `yaml:"Redis"`
}

type RedisConfig struct {
	Host      []string `yaml:"Host"`
	Password  string   `yaml:"Password"`
	DB        int      `yaml:"DB"`
	MaxIdle   int      `yaml:"MaxIdle"`
	MaxActive int      `yaml:"MaxActive"`
}

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err = viper.Unmarshal(&Instance)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	fullSheetRange(&Instance)

	// auto reload config if there is any changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.Unmarshal(&Instance)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}
		fullSheetRange(&Instance)
	})
}

func fullSheetRange(config *Config) {
	config.MaxFetchRange = fmt.Sprintf("%s!%s", config.SheetName, config.MaxFetchRange)
}
