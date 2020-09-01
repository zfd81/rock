package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Banner  string `mapstructure:"banner"`
	Meta    Meta   `mapstructure:"meta"`
	Etcd    Etcd   `mapstructure:"etcd"`
}

type Meta struct {
	PathSeparator    string `mapstructure:"pathSeparator"` // 路径分隔符（分隔路径元素）
	NameSeparator    string `mapstructure:"nameSeparator"` // 名字分隔符（分隔对象全名）
	RootDirectory    string `mapstructure:"rootDirectory"`
	ServiceDirectory string `mapstructure:"rootDirectory"`
}

type Etcd struct {
	Endpoints      []string `mapstructure:"endpoints"`
	DialTimeout    int      `mapstructure:"dialTimeout"`
	RequestTimeout int      `mapstructure:"requestTimeout"`
}

const (
	ConfigName = "parrot-site"
	ConfigPath = "."
	ConfigType = "yaml"

	banner1 = `
	 ____   __    ____  ____  _____  ____ 
	(  _ \ /__\  (  _ \(  _ \(  _  )(_  _)
	 )___//(__)\  )   / )   / )(_)(   )(  
	(__) (__)(__)(_)\_)(_)\_)(_____) (__) 

	`

	banner2 = `

	########     ###    ########  ########   #######  ######## 
	##     ##   ## ##   ##     ## ##     ## ##     ##    ##    
	##     ##  ##   ##  ##     ## ##     ## ##     ##    ##    
	########  ##     ## ########  ########  ##     ##    ##    
	##        ######### ##   ##   ##   ##   ##     ##    ##    
	##        ##     ## ##    ##  ##    ##  ##     ##    ##    
	##        ##     ## ##     ## ##     ##  #######     ## 
	`
)

var defaultConf = Config{
	Name:    "Parrot",
	Version: "1.0.0",
	Banner:  banner1,
	Meta: Meta{
		PathSeparator:    "/",
		NameSeparator:    ".",
		RootDirectory:    "/parrot",
		ServiceDirectory: "/serv",
	},
	Etcd: Etcd{
		Endpoints:      []string{"127.0.0.1:2379"},
		DialTimeout:    5,
		RequestTimeout: 5,
	},
}

var globalConf = defaultConf

func init() {
	viper.SetConfigName(ConfigName)
	viper.AddConfigPath(ConfigPath)
	viper.SetConfigType(ConfigType)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Errorf("Fatal error config file: %s \n", err))
	} else {
		err = viper.Unmarshal(&globalConf)
		if err != nil {
			panic(fmt.Errorf("Fatal error when reading %s config, unable to decode into struct, %v", ConfigName, err))
		}
	}
}

func GetConfig() *Config {
	return &globalConf
}
