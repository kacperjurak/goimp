package goimp

import "gonum.org/v1/gonum/optimize"

type Weighting int

const (
	MODULUS Weighting = iota
	UNITY
)

type Result struct {
	Params          []float64       `json:"params"`
	ChiSq           float64         `json:"chiSq"`
	MajorIterations int             `json:"majorIterations"`
	FuncEvaluations int             `json:"funcEvaluations"`
	Runtime         float64         `json:"runtime"`
	Status          optimize.Status `json:"status"`
}

type Core interface {
	CircuitImpedance(code string, freqs []float64, values []float64) []complex128
	Solve(freqs []float64, impData [][2]float64) (Result, error)
}
