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
