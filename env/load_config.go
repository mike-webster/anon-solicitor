package env

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	yaml "gopkg.in/yaml.v1"
)

var path = "app.yaml"
var curEnv *Environment
var target = "development"

type Environment struct {
	ConnectionString string `yaml:"connection_string"`
	SMTPHost         string `yaml:"smtp_host"`
	SMTPPort         int    `yaml:"smtp_port"`
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
}

func Target() string {
	return target
}

// Config always returns the configuration settings
// for currently configured environment.
func Config() *Environment {
	if curEnv == nil {
		setTarget()
		curEnv = loadAppConfig()
	}

	log.Printf("-- [LOADED CONFIG]\n\t~ SMTPHost: %v\n\t~ SMTPPort: %v\n\t~ Port: %v", curEnv.SMTPHost, curEnv.SMTPPort, curEnv.Port)

	return curEnv
}

func setTarget() {
	target := os.Getenv("GO_ENV")
	switch target {
	case "development", "production", "uat", "test":
	default:
		panic(fmt.Errorf("Invalid target: %v", target))
	}
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

	res := envs[Target()]
	envconfig.MustProcess("", &res)

	return &res
}
