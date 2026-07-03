package bootstrap

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/logging"
)

type Config struct {
	Constants *Constants
	Env       *Env
}

func Run() *Config {
	env := NewEnv()
	logging.Setup(env.Server.Mode)

	return &Config{
		Constants: NewConstants(),
		Env:       env,
	}
}
