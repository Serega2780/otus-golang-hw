package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger                 *LoggerConf
	IsInMemory             bool `split_words:"true" default:"true" yaml:"isInMemory"`
	DB                     *DBConf
	HTTP                   *HTTPServerConfig
	GRPC                   *GRPCServerConfig
	AMQP                   *AMQPConfig
	SchedulerPeriodSeconds int `split_words:"true" default:"2" yaml:"schedulerPeriodSeconds"`
}

type LoggerConf struct {
	Level        string `default:"info" yaml:"level"`
	Format       string `default:"text" yaml:"format"`
	File         string `default:"calendar.log" yaml:"file"`
	LogToFile    bool   `split_words:"true" default:"false" yaml:"logToFile"`
	LogToConsole bool   `split_words:"true" default:"true" yaml:"logToConsole"`
}

type DBConf struct {
	User      string `default:"user" yaml:"user"`
	Password  string `default:"dummy" yaml:"password"`
	Dbname    string `default:"calendar" yaml:"dbname"`
	Host      string `default:"localhost" yaml:"host"`
	Port      string `default:"5432" yaml:"port"`
	Migration string `default:"migrations" yaml:"migration"`
}

type HTTPServerConfig struct {
	IP   string `default:"0.0.0.0" yaml:"ip"`
	Port string `default:"8585" yaml:"port"`
}

type GRPCServerConfig struct {
	IP   string `default:"0.0.0.0" yaml:"ip"`
	Port string `default:"6565" yaml:"port"`
}

type AMQPConfig struct {
	IP        string `default:"0.0.0.0" yaml:"ip"`
	Port      string `default:"5672" yaml:"port"`
	QueueName string `split_words:"true" default:"event_notify" yaml:"queueName"`
}

func New() *Config {
	return &Config{
		IsInMemory:             true,
		SchedulerPeriodSeconds: 30,
		Logger:                 &LoggerConf{Level: "info", Format: "text", LogToFile: false, LogToConsole: true},
		DB:                     &DBConf{},
		HTTP:                   &HTTPServerConfig{IP: "0.0.0.0", Port: "8585"},
		GRPC:                   &GRPCServerConfig{IP: "0.0.0.0", Port: "6565"},
		AMQP:                   &AMQPConfig{IP: "0.0.0.0", Port: "5672"},
	}
}

func ReadConfig(configFile string) *Config {
	conf := &Config{}
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Could not open a config file %v\n", err)
	} else {
		err = yaml.Unmarshal(yamlFile, conf)
		if err != nil {
			fmt.Printf("Could not unmarshal a config file %v\n", err)
		}
	}
	if conf.Logger == nil {
		err = envconfig.Process("", conf)
		if err != nil {
			fmt.Printf("Process ENV variables %v\n", err)
			fmt.Println("Will use default config")
			conf = New()
		}
	}
	return conf
}
