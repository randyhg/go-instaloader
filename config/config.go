package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go-instaloader/utils/rlog"
	"log"

	"github.com/spf13/viper"
)

var Instance Config

type Config struct {
	Host       string `yaml:"Host"`
	SocketPort string `yaml:"SocketPort"`
	Env        string `yaml:"Env"`

	PyInstaloaderDomain string `yaml:"PyInstaloaderDomain"`

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
	ServiceKeyPath string `yaml:"ServiceKeyPath"`
	CredentialPath string `yaml:"CredentialPath"`
	TokenPath      string `yaml:"TokenPath"`

	// db
	ShowSql      bool        `yaml:"ShowSql"`
	MySqlUrl     string      `yaml:"MySqlUrl"`
	MySqlMaxIdle int         `yaml:"MySqlMaxIdle"`
	MySqlMaxOpen int         `yaml:"MySqlMaxOpen"`
	Redis        RedisConfig `yaml:"Redis"`

	// log
	LogDir      string `yaml:"LogDir"`
	LogFileName string `yaml:"LogFileName"`
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
	WriteLog()
	fullSheetRange(&Instance)

	// auto reload config if there is any changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.Unmarshal(&Instance)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}
		WriteLog()
		fullSheetRange(&Instance)
	})

}

func fullSheetRange(config *Config) {
	config.MaxFetchRange = fmt.Sprintf("%s!%s", config.SheetName, config.MaxFetchRange)
}

func WriteLog() {
	if Instance.Env != "" {
		l := rlog.New(
			Instance.LogDir,
			Instance.LogFileName,
			Instance.Env,
			true,
			log.LstdFlags|log.Lshortfile)
		rlog.Export(l)
	}
}
