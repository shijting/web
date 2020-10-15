package inits

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)
type AppConfig struct {
	Name      string `mapstructure:"name"`
	Version   string `mapstructure:"version"`
	*DatabaseConfig `mapstructure:"database"`
	*RedisConfig `mapstructure:"redis"`
}
type DatabaseConfig struct {
	Addr     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Passwrod string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
}

func InitConfig(filePath string) (err error) {

	viper.SetConfigFile(filePath)

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return
}