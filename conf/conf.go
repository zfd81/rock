package conf

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Name         string        `mapstructure:"name"`
	Version      string        `mapstructure:"version"`
	Banner       string        `mapstructure:"banner"`
	Port         int           `mapstructure:"port"`
	ServiceName  string        `mapstructure:"service-name"`
	Directory    string        `mapstructure:"directory"`
	ReadTimeout  time.Duration `mapstructure:"read-timeout"`
	WriteTimeout time.Duration `mapstructure:"write-timeout"`
	APIServer    APIServer     `mapstructure:"api-server"`
	Meta         Meta          `mapstructure:"meta"`
	Etcd         Etcd          `mapstructure:"etcd"`
	Cluster      Cluster       `mapstructure:"cluster"`
}

type APIServer struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read-timeout"`
	WriteTimeout time.Duration `mapstructure:"write-timeout"`
}

type Meta struct {
	PathSeparator string `mapstructure:"path-separator"` // 路径分隔符（分隔路径元素）
	NameSeparator string `mapstructure:"name-separator"` // 名字分隔符（分隔对象全名）
}

type Etcd struct {
	Endpoints      []string `mapstructure:"endpoints"`
	DialTimeout    int      `mapstructure:"dial-timeout"`
	RequestTimeout int      `mapstructure:"request-timeout"`
}

type Cluster struct {
	HeartbeatInterval        int `mapstructure:"heartbeat-interval"`
	HeartbeatRecheckInterval int `mapstructure:"heartbeat-recheck-interval"`
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
	Name:         "Parrot",
	Version:      "1.0.0",
	Banner:       banner1,
	Port:         8081,
	ServiceName:  "/",
	Directory:    "/parrot",
	ReadTimeout:  10,
	WriteTimeout: 15,
	APIServer: APIServer{
		Port:         8143,
		ReadTimeout:  5,
		WriteTimeout: 10,
	},
	Meta: Meta{
		PathSeparator: "/",
		NameSeparator: ".",
	},
	Etcd: Etcd{
		Endpoints:      []string{"127.0.0.1:2379"},
		DialTimeout:    5,
		RequestTimeout: 5,
	},
	Cluster: Cluster{
		HeartbeatInterval:        9,
		HeartbeatRecheckInterval: 5,
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
