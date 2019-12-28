package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kacperjurak/goimp"
	"github.com/kacperjurak/goimpcore"
	"runtime"
	"sync"
	"time"
)

func HandleRequest(r goimp.Request) (goimp.Response, error) {
	var response goimp.Response

	taskLen := len(r.Tasks)
	var resultsIndexed = make([]goimp.Result, taskLen)
	c := make(chan goimp.ResultIndexed, taskLen)
	var wg sync.WaitGroup

	start := time.Now()

	for i, t := range r.Tasks {
		wg.Add(1)
		go solve(r.Config, t, &wg, c, i)
	}

	wg.Wait()
	close(c)

	response.Runtime = float64(time.Since(start) / 1000)
	response.SpectrumsNo = taskLen
	response.Code = r.Code
	response.MaxProcs = runtime.GOMAXPROCS(0)
	for result := range c {
		resultsIndexed[result.Index] = result.Result
	}
	response.Data = resultsIndexed

	return response, nil
}

func solve(config goimp.Config, task goimp.Task, wg *sync.WaitGroup, c chan<- goimp.ResultIndexed, index int) {
	defer wg.Done()
	var (
		freqs   = task.Freqs
		impData = task.ImpData
	)

	freqs = freqs[task.CutLow : len(freqs)-int(task.CutHigh)]
	impData = impData[task.CutLow : len(impData)-int(task.CutHigh)]

	var s goimp.Solver = goimpcore.NewSolver(config.Code, task.InitValues, goimpcore.MODULUS)

	r, err := s.Solve(freqs, impData)
	if err != nil {
		panic(err)
	}

	c <- goimp.ResultIndexed{
		Index:  index,
		Result: r,
	}
}

func main() {
	lambda.Start(HandleRequest)
}
