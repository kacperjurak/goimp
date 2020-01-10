package goimp

const (
	OK Status = iota
	NOK
)

type Status int

type Task struct {
	Freqs      []float64    `json:"freqs"`
	ImpData    [][2]float64 `json:"impData"`
	InitValues []float64    `json:"initValues"`
	CutLow     uint         `json:"cutLow"`
	CutHigh    uint         `json:"cutHigh"`
}

type TaskIndexed struct {
	Index int
	Task  Task
}

type Config struct {
	Code    string `json:"code"`
	CutLow  uint   `json:"cutLow"`
	CutHigh uint   `json:"cutHigh"`
}

type Request struct {
	Config
	Tasks []Task `json:"tasks"`
}

type Result struct {
	Params  []float64   `json:"params"`
	Min     float64     `json:"min"`
	MinUnit string      `json:"minUnit"`
	Payload interface{} `json:"payload"`
	Runtime float64     `json:"runtime"`
	Status  Status      `json:"status"`
}

type ResultIndexed struct {
	Index  int
	Result Result
}

type Response struct {
	Code        string   `json:"code"`
	SpectrumsNo int      `json:"spectrumsNo"`
	Runtime     float64  `json:"runtime"`
	MaxProcs    int      `json:"maxProcs"`
	Data        []Result `json:"data"`
}

type Solver interface {
	Solve(freqs []float64, impData [][2]float64) (Result, error)
}
