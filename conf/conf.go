package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"readingo/constant"
)

var (
	Server        ServerConf
	Redis         []RedisConf
	RedisClusters []RedisClusterConf
	Auth          AuthConf
)

type config struct {
	Server        ServerConf         `yaml:"server"`
	Redis         []RedisConf        `yaml:"redis"`
	RedisClusters []RedisClusterConf `yaml:"redisc"`
	Auth          AuthConf           `yaml:"auth"`
}

type RedisClusterConf struct {
	HOST           []string `yaml:"host"`
	PASSWORD       string   `yaml:"password"`
	NAME           string   `yaml:"name"`
	MaxIdle        int      `yaml:"maxIdle"`
	MaxActive      int      `yaml:"maxActive"`
	IdleTimeout    string   `yaml:"idleTimeout"`
	ConnectTimeout string   `yaml:"connectTimeout"`
}

type ServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type UserConf struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Role     string `yaml:"role"`
}

type RedisConf struct {
	HOST           string `yaml:"host"`
	PASSWORD       string `yaml:"password"`
	NAME           string `yaml:"name"`
	MaxIdle        int    `yaml:"maxIdle"`
	MaxActive      int    `yaml:"maxActive"`
	IdleTimeout    string `yaml:"idleTimeout"`
	ConnectTimeout string `yaml:"connectTimeout"`
}

type AuthConf struct {
	Anonymous   bool       `yaml:"anonymous"`
	DefaultRole string     `yaml:"defaultRole"`
	Users       []UserConf `yaml:"users"`
}

func init() {
	configFile, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		panic("conf.yml not found, please at least initialize one")
	}

	var conf config
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		panic("Cannot unmarshal conf.yml, error: " + err.Error())
	}
	if conf.Auth.DefaultRole == "" {
		conf.Auth.DefaultRole = constant.RoleReadOnly
	}

	Server = conf.Server
	Redis = conf.Redis
	RedisClusters = conf.RedisClusters
	Auth = conf.Auth
}
