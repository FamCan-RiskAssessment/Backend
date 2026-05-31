package calcdto

import formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"

type CalcEnumResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ModelResponse struct {
	Name        string  `json:"name"`
	Probability float64 `json:"probability"`
}

type Premm5Response struct {
	GeneProbs map[string]float64 `json:"gene_probs"`
	PAny      float64            `json:"p_any"`
	PNone     float64            `json:"p_none"`
}

type BCRAResponse struct {
	AbsRisk    float64 `json:"AbsRisk"`
	AbsRiskAvg float64 `json:"AbsRisk_Avg"`
	RRStar1    float64 `json:"RR_Star1"`
	RRStar2    float64 `json:"RR_Star2"`
	ProjIntvl  float64 `json:"Proj_Intvl"`
}

type GailResponse struct {
	AbsoluteRisk float64  `json:"absolute_risk"`
	RelativeRisk *float64 `json:"relative_risk"`
}

type PLCOResponse struct {
	PLCOM20126YrRisk     float64            `json:"plcom2012_6yr_risk"`
	PLCOM20123YrRisk     float64            `json:"plcom2012_3yr_risk"`
	PLCOM2012RiskPercent float64            `json:"plcom2012_risk_percent"`
	PLCO2012Results      map[string]float64 `json:"plco2012results,omitempty"`
}

type CalcResultsBundle struct {
	Premm5 *Premm5Response `json:"premm5,omitempty"`
	BCRA   *BCRAResponse   `json:"bcra,omitempty"`
	Gail   *GailResponse   `json:"gail,omitempty"`
	PLCO   *PLCOResponse   `json:"plco,omitempty"`
}

type CalcBrowseItemResponse struct {
	Form    formdto.BasicFormResponse `json:"form"`
	Results CalcResultsBundle         `json:"results"`
}
