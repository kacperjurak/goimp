package main

import (
	"encoding/json"
	"flag"
	"github.com/kacperjurak/goimp"
	"github.com/kacperjurak/goimpcore"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var verbose bool

func main() {
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.Parse()
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var response goimp.Response

	switch r.Method {
	case "POST":
		d := json.NewDecoder(r.Body)
		p := &goimp.Request{}
		err := d.Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		taskLen := len(p.Tasks)
		var resultsIndexed = make([]goimp.Result, taskLen)
		c := make(chan goimp.ResultIndexed, taskLen)
		var wg sync.WaitGroup

		start := time.Now()

		for i, t := range p.Tasks {
			wg.Add(1)
			go solve(p.Config, t, &wg, c, i)
		}

		wg.Wait()
		close(c)

		response.Runtime = float64(time.Since(start) / 1000)
		response.SpectrumsNo = taskLen
		response.Code = p.Code
		response.MaxProcs = runtime.GOMAXPROCS(0)
		for result := range c {
			resultsIndexed[result.Index] = result.Result
		}
		response.Results = resultsIndexed
		j, _ := json.Marshal(response)

		if _, err := w.Write(j); err != nil {
			panic(err)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func solve(config goimp.Config, task goimp.Task, wg *sync.WaitGroup, c chan<- goimp.ResultIndexed, index int) {
	defer wg.Done()
	var (
		freqs   = task.Freqs
		impData = task.ImpData
	)

	freqs = freqs[task.CutLow : len(freqs)-int(task.CutHigh)]
	impData = impData[task.CutLow : len(impData)-int(task.CutHigh)]

	s := goimpcore.NewSolver(config.Code)
	s.InitValues = task.InitValues
	s.SmartMode = "eis"

	r, err := s.Solve(freqs, impData)
	if err != nil {
		panic(err)
	}

	if verbose {
		log.Println(index, r)
	}

	c <- goimp.ResultIndexed{
		Index:  index,
		Result: r,
	}
}
