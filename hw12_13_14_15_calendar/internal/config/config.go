package config

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     *LoggerConf
	IsInMemory bool `split_words:"true" default:"true" yaml:"isInMemory"`
	DB         *DBConf
	HTTP       *HTTPServerConfig
	// TODO
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

func New() *Config {
	return &Config{
		IsInMemory: true,
		Logger:     &LoggerConf{Level: "info", Format: "text", LogToFile: false, LogToConsole: true},
		DB:         &DBConf{},
		HTTP:       &HTTPServerConfig{IP: "0.0.0.0", Port: "8585"},
	}
}
