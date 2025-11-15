package bootstrap

import "sync"

type Config struct {
	Constants *Constants
	Env       *Env
}

var runOnce sync.Once

func Run() *Config {
	runOnce.Do(func() {
		projectConfig = &Config{
			Constants: NewConstants(),
			Env:       NewEnv(),
		}
	})
	return projectConfig
}

var projectConfig *Config