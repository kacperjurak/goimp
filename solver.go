package goimp

const (
	OK Status = iota
	NOK
)

type Status int

type Task struct {
	Code       string       `json:"code"`
	Freqs      []float64    `json:"freqs"`
	ImpData    [][2]float64 `json:"impData"`
	InitValues []float64    `json:"initValues"`
}

type Request struct {
	Index int  `json:"index"`
	Task  Task `json:"task"`
}

type Result struct {
	Code       string      `json:"code"`
	InitValues []float64   `json:"initValues"`
	Params     []float64   `json:"params"`
	Min        float64     `json:"min"`
	MinUnit    string      `json:"minUnit"`
	Payload    interface{} `json:"payload"`
	Runtime    float64     `json:"runtime"`
	Status     Status      `json:"status"`
}

type Response struct {
	Index  int    `json:"index"`
	Result Result `json:"result"`
}

type Solver interface {
	Solve(freqs []float64, impData [][2]float64) (Result, error)
}
