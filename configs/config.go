package configs

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type HttpConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type Kubernetes struct {
	Conf string `yaml:"conf"`
}

type Config struct {
	Http       HttpConfig `yaml:"HttpService"`
	LogFile    string     `yaml:"LogFile"`
	Kubernetes Kubernetes `json:"Kubernetes"`
	ConfigPath string
}

var (
	conf       *Config
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configPath", "lunara-k8s.yaml", "configuration file path")
}

func Init() (err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile(configPath); err != nil {
		return errors.New(fmt.Sprintf("ioutil.ReadFile err:%v", err))
	}
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return errors.New(fmt.Sprintf("yaml.Unmarshal err:%v", err))
	}
	conf.ConfigPath = configPath
	return nil
}

func GetConfig() *Config {
	log.Println(conf)
	return conf
}
