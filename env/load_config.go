package env

import (
	"io/ioutil"
	"log"
	"path/filepath"
	_ "path/filepath"

	"github.com/kelseyhightower/envconfig"
	yaml "gopkg.in/yaml.v1"
)

var path = "app.yaml"
var curEnv *Environment

type Environment struct {
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort string `yaml:"smtp_port"`
	Port     string `yaml:"port"`
}

// Config always returns the configuration settings
// for currently configured environment.
func Config() *Environment {
	if curEnv == nil {
		curEnv = loadAppConfig()
	}

	log.Printf("-- [LOADED CONFIG]\n\t~ SMTPHost: %v\n\t~ SMTPPort: %v\n\t~ Port: %v", curEnv.SMTPHost, curEnv.SMTPPort, curEnv.Port)

	return curEnv
}

func loadAppConfig() *Environment {
	fp, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	f, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	var envs map[string]Environment
	if err = yaml.Unmarshal(f, &envs); err != nil {
		panic(err)
	}

	res := envs["development"]
	envconfig.MustProcess("", &res)

	return &res
}
