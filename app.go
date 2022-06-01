package devproxy

import (
	"net/url"

	"azugo.io/devproxy/spa"

	"azugo.io/azugo"
	"azugo.io/azugo/server"
	"go.uber.org/zap"
)

type App struct {
	*azugo.App
	config *Config

	spa spa.SpaDevProxy

	routesAdded bool
}

func NewApp() (*App, error) {
	a, err := server.New(nil, server.ServerOptions{
		AppName: "Azugo DevProxy",
	})
	if err != nil {
		return nil, err
	}
	return &App{
		App: a,
		config: &Config{
			Proxy: make(map[string]*ConfigProxy),
		},
	}, nil
}

func (a *App) Start() error {
	if !a.routesAdded {
		a.routesAdded = true
		for _, c := range a.config.Proxy {
			v := c
			u, err := url.Parse(v.Upstream)
			if err != nil {
				a.Log().With(zap.Error(err)).Fatal("failed to parse upstream URL")
				return err
			}
			a.Proxy(v.Path, azugo.ProxyUpstream(u))
		}

		var err error
		switch a.config.Spa.Type {
		case "vue":
			a.spa, err = spa.NewVueDevProxy(a.config.Spa.Vue)
			if err != nil {
				return err
			}
		}

		// Forward to frontend
		a.Proxy("/", azugo.ProxyUpstream(a.spa.DevServerURL()))
	}

	if err := a.spa.Start(a.BackgroundContext()); err != nil {
		return err
	}

	return a.App.Start()
}

func (a *App) Stop() {
	if err := a.spa.Stop(); err != nil {
		// Ignore signal killed error
		if err.Error() != "signal: killed" {
			a.Log().With(zap.Error(err)).Warn("failed to stop SPA")
		}
	}
	a.App.Stop()
}
