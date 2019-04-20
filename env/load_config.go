package env

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	SMTPUser         string `yaml:"smtp_user"`
	SMTPPass         string `envconfig:"SMTP_PASS"`
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	Secret           string `envconfig:"ANON_SOLICITOR_SECRET"`
	AppName          string `envconfig:"APP_NAME"`
	ShouldSendEmails bool   `envconfig:"SEND_EMAILS"`
}

func (e *Environment) ToString() string {
	ret := "-- [LOADED CONFIG]\n\t"
	ret += "~ Connection String: %v\n\t"
	ret += "~ Secret: %v\n\t"
	ret += "~ SMTPHost: %v\n\t"
	ret += "~ SMTPPort: %v\n\t"
	ret += "~ Host: %v\n\t"
	ret += "~ Port: %v\n\t"
	ret += "~ AppName: %v\n\t"

	return fmt.Sprintf(ret,
		e.ConnectionString,
		e.Secret,
		e.SMTPHost,
		e.SMTPPort,
		e.Host,
		e.Port,
		e.AppName)
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
		log.Printf(curEnv.ToString())
	}

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
	// So.. this is my shitty work around so we can
	// call env.Config() from test files in other packages
	// and still load the config correctly.
	var fp string
	wd, _ := os.Getwd()
	appName := os.Getenv("APP_NAME")
	if strings.HasSuffix(wd, appName) {
		tp, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		fp = tp
	} else {
		cut := strings.LastIndex(wd, "/")
		fp = wd[:cut+1] + path
	}

	f, err := ioutil.ReadFile("" + fp)
	if err != nil {
		log.Print(err)
		return nil
	}

	var envs map[string]Environment
	if err = yaml.Unmarshal(f, &envs); err != nil {
		panic(err)
	}

	res := envs[Target()]
	envconfig.MustProcess("", &res)

	return &res
}
