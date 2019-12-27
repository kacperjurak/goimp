package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kacperjurak/goimpcore"
	"sync"
)

type Task struct {
	Code       string       `json:"code"`
	Freqs      []float64    `json:"freqs"`
	ImpData    [][2]float64 `json:"impData"`
	InitValues []float64    `json:"initValues"`
	CutLow     uint         `json:"cutLow"`
	CutHigh    uint         `json:"cutHigh"`
}

type Request struct {
	Tasks []Task `json:"tasks"`
}

type Result struct {
	Index int `json:"-"`
	goimpcore.Result
}

type Response struct {
	Data []Result `json:"data"`
}

func HandleRequest(_ context.Context, r Request) (Response, error) {
	var response Response

	results := make(chan Result, len(r.Tasks))
	var wg sync.WaitGroup

	for i, t := range r.Tasks {
		wg.Add(1)
		go solve(t, &wg, results, i)
	}

	wg.Wait()
	close(results)
	for result := range results {
		response.Data = append(response.Data, result)
	}

	return response, nil
}

func solve(task Task, wg *sync.WaitGroup, c chan<- Result, index int) {
	defer wg.Done()
	var (
		freqs   = task.Freqs
		impData = task.ImpData
	)

	freqs = freqs[task.CutLow : len(freqs)-int(task.CutHigh)]
	impData = impData[task.CutLow : len(impData)-int(task.CutHigh)]

	s := goimpcore.NewSolver(task.Code, task.InitValues, goimpcore.MODULUS)
	r, err := s.Solve(freqs, impData)
	if err != nil {
		panic(err)
	}
	c <- Result{
		Index:  index,
		Result: r,
	}
}

func main() {
	lambda.Start(HandleRequest)
}
