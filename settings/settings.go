package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

// 实例化对象：对应yaml 文件
var Conf = new(AppConfig)

type AppConfig struct {
	*LogConfig   `mapstructure:"log"`
	*HttpConfig  `mapstructure:"http"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
	*EnvConfig   `mapstructure:"env"`
}

type EnvConfig struct {
	Crontab bool `mapstructure:"crontab"`
}

type HttpConfig struct {
	BaseUrl   string `mapstructure:"base_url"`
	Ip        string `mapstructure:"ip"`
	AdminPort string `mapstructure:"admin_port"`
	GinMode   string `mapstructure:"gin_mode"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	*Options `mapstructure:"options"`
}
type Options struct {
	Loc       string `mapstructure:"loc"`
	Charset   string `mapstructure:"charset"`
	Collation string `mapstructure:"collation"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Protocol string `mapstructure:"protocol"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

func InitConfig(configPaths string) {
	str, _ := os.Getwd()
	viper.SetConfigFile(str + configPaths)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("ReadInConfig failed,err:%v", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		//重载
		viper.Unmarshal(&Conf)
	})
	// 将配置赋值给全局变量
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("Unmarshal failed,err:%v", err))
	}
}
