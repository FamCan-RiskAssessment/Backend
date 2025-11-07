package calcdto

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
