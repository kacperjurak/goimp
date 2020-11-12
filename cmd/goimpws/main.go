package main

import (
	"encoding/json"
	"flag"
	"github.com/kacperjurak/goimp"
	"github.com/kacperjurak/goimpcore"
	"log"
	"net/http"
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

		//start := time.Now()

		result := solve(p.Task, p.Index)

		//response.Result.Runtime = float64(time.Since(start) / 1000)
		//response.Result.Code = p.Task.Code
		//response.Result.MaxProcs = runtime.GOMAXPROCS(0)
		response.Index = p.Index
		response.Result = result
		j, _ := json.Marshal(response)

		if _, err := w.Write(j); err != nil {
			panic(err)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func solve(task goimp.Task, index int) goimp.Result {
	var (
		freqs   = task.Freqs
		impData = task.ImpData
	)

	s := goimpcore.NewSolver(task.Code)
	s.InitValues = task.InitValues
	s.SmartMode = "eis"

	r, err := s.Solve(freqs, impData)
	if err != nil {
		panic(err)
	}

	if verbose {
		log.Println(index, r)
	}

	return r
}
