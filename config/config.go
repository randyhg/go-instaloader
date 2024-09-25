package config

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go-instaloader/utils/rlog"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var Instance Config

type Config struct {
	Host string `yaml:"Host"`
	Env  string `yaml:"Env"`

	PyInstaloaderDomain string `yaml:"PyInstaloaderDomain"`

	// spreadsheet
	SpreadSheetId   string `yaml:"SpreadSheetId"`
	SheetName       string `yaml:"SheetName"`
	MaxFetchRange   string `yaml:"MaxFetchRange"`
	StatusColumn    string `yaml:"StatusColumn"`
	RemarkColumn    string `yaml:"RemarkColumn"`
	ConfigSheetName string `yaml:"ConfigSheetName"`
	ConfigCellRange string `yaml:"ConfigCellRange"`

	// queue
	DelayWhenNoJobInSeconds   int `yaml:"DelayWhenNoJobInSeconds"`
	DelayWhenErrorInSeconds   int `yaml:"DelayWhenErrorInSeconds"`
	DelayWhenJobDoneInSeconds int `yaml:"DelayWhenJobDoneInSeconds"`

	// google api cred
	ServiceKeys    []string `yaml:"ServiceKeys"`
	CredentialPath string   `yaml:"CredentialPath"`
	TokenPath      string   `yaml:"TokenPath"`

	// db
	ShowSql      bool        `yaml:"ShowSql"`
	MySqlUrl     string      `yaml:"MySqlUrl"`
	MySqlMaxIdle int         `yaml:"MySqlMaxIdle"`
	MySqlMaxOpen int         `yaml:"MySqlMaxOpen"`
	Redis        RedisConfig `yaml:"Redis"`

	// log
	LogDir      string `yaml:"LogDir"`
	LogFileName string `yaml:"LogFileName"`

	// tg bot
	TeleBotToken string `yaml:"TeleBotToken"`
	TeleGroupId  string `yaml:"TeleGroupId"`
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
		rlog.Fatalf("Error reading config file, %s", err)
	}
	err = viper.Unmarshal(&Instance)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	configInitFunc()

	// auto reload config if there is any changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.Unmarshal(&Instance)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}
		configInitFunc()
	})

}

func configInitFunc() {
	WriteLog()
	fullSheetRange(&Instance)
	getServiceKeyEmail(Instance.ServiceKeys)
}

func fullSheetRange(config *Config) {
	config.MaxFetchRange = fmt.Sprintf("%s!%s", config.SheetName, config.MaxFetchRange)
}

func getServiceKeyEmail(files []string) {
	var emails []string
	for _, file := range files {
		fileByte, err := os.ReadFile(file)
		if err != nil {
			rlog.Errorf("error reading file %s: %v", file, err)
			continue
		}

		var serviceKey map[string]interface{}
		if err = json.Unmarshal(fileByte, &serviceKey); err != nil {
			rlog.Errorf("error unmarshalling file %s: %v", file, err)
			continue
		}
		email, ok := serviceKey["client_email"].(string)
		if ok {
			emails = append(emails, email)
		}
	}
	emailList := strings.Join(emails, "\n")
	fmt.Printf("============================>\nplease grant access to the spreadsheet for these emails:\n%s\n============================>\n", emailList)
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
