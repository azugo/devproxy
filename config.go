package devproxy

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"azugo.io/devproxy/spa"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type ConfigSpa struct {
	Type string                  `yaml:"type"`
	Vue  *spa.VueDevProxyOptions `yaml:"vue"`
}

type ConfigProxy struct {
	Path     string `yaml:"path"`
	Upstream string `yaml:"upstream"`
}

type Config struct {
	Spa   *ConfigSpa              `yaml:"spa"`
	Proxy map[string]*ConfigProxy `yaml:"proxy"`
}

func (a *App) LoadConfig() error {
	f, err := os.Open(".devproxy.yml")
	if err == nil && f != nil {
		defer f.Close()
		buf, err := io.ReadAll(f)
		if err != nil {
			a.Log().With(zap.Error(err)).Warn("failed to read .devproxy.yml file")
		} else if err = yaml.Unmarshal(buf, a.config); err != nil {
			a.Log().With(zap.Error(err)).Warn("failed to parse .devproxy.yml file")
		}
	}

	// Allow to override upstream URL via environment variable
	for name, proxy := range a.config.Proxy {
		if e := os.Getenv(fmt.Sprintf("PROXY_%s_UPSTREAM", strings.ToUpper(name))); len(e) > 0 {
			proxy.Upstream = e
		}
	}

	if a.config.Spa == nil {
		a.config.Spa = &ConfigSpa{}
	}
	if a.config.Spa.Type == "" {
		a.config.Spa.Type = "vue"
	}
	if a.config.Spa.Type == "vue" {
		if a.config.Spa.Vue == nil {
			a.config.Spa.Vue = &spa.VueDevProxyOptions{
				RunnerType: spa.RunnerTypeYarn,
				Dir:        ".",
			}
		}
	} else {
		return errors.New("unknown spa type")
	}

	return nil
}
