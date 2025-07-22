package wire

import "github.com/FamCan-RiskAssessment/Backend/bootstrap"

type Application struct {
}

func InitializeApplication(config *bootstrap.Config) (*Application, error) {
	return &Application{}, nil
}
