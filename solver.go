package goimp

type Result struct {
	Params          []float64 `json:"params"`
	ChiSq           float64   `json:"chiSq"`
	MajorIterations int       `json:"majorIterations"`
	FuncEvaluations int       `json:"funcEvaluations"`
	Runtime         float64   `json:"runtime"`
	Status          string    `json:"status"`
}

type Solver interface {
	Solve(freqs []float64, impData [][2]float64) (Result, error)
}
