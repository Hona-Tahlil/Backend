package bootstrap

type Config struct {
	Constants *Constants
	Env       *Env
}

func Run() *Config {
	ProjectConfig = &Config{
		Constants: NewConstants(),
		Env:       NewEnv(),
	}
	return ProjectConfig
}

var ProjectConfig *Config
