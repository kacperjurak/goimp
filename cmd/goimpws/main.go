package main

import (
	"encoding/json"
	"flag"
	"github.com/kacperjurak/goimp"
	"github.com/kacperjurak/goimpcore"
	"log"
	"net/http"
	"sync"
	"time"
)

type Task struct {
	Freqs      []float64    `json:"freqs"`
	ImpData    [][2]float64 `json:"impData"`
	InitValues []float64    `json:"initValues"`
	CutLow     uint         `json:"cutLow"`
	CutHigh    uint         `json:"cutHigh"`
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
	Index int `json:"-"`
	goimp.Result
}

type Response struct {
	Code    string   `json:"code"`
	Runtime float64  `json:"runtime"`
	Data    []Result `json:"data"`
}

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
	var response Response

	switch r.Method {
	case "POST":
		d := json.NewDecoder(r.Body)
		p := &Request{}
		err := d.Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		results := make(chan Result, len(p.Tasks))
		var wg sync.WaitGroup

		start := time.Now()

		for i, t := range p.Tasks {
			wg.Add(1)
			go solve(p.Config, t, &wg, results, i)
		}

		wg.Wait()
		close(results)

		response.Runtime = float64(time.Since(start) / 1000)

		for result := range results {
			response.Data = append(response.Data, result)
		}
		j, _ := json.Marshal(response)

		if _, err := w.Write(j); err != nil {
			panic(err)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func solve(config Config, task Task, wg *sync.WaitGroup, c chan<- Result, index int) {
	defer wg.Done()
	var (
		freqs   = task.Freqs
		impData = task.ImpData
	)

	freqs = freqs[task.CutLow : len(freqs)-int(task.CutHigh)]
	impData = impData[task.CutLow : len(impData)-int(task.CutHigh)]

	log.Println(config.Code)
	var s goimp.Solver
	s = goimpcore.NewSolver(config.Code, task.InitValues, goimpcore.MODULUS)

	r, err := s.Solve(freqs, impData)
	if err != nil {
		panic(err)
	}

	if verbose {
		log.Println(index, r)
	}

	c <- Result{
		Index:  index,
		Result: r,
	}
}
