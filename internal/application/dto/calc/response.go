package calcdto

type CalcEnumResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ModelResponse struct {
	Name        string  `json:"name"`
	Probability float64 `json:"probability"`
}
