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

		result := solve(p.Task)

		response.Index = p.Index
		response.Result = result
		response.Result.InitValues = p.Task.InitValues
		j, _ := json.Marshal(response)

		if verbose {
			log.Println(result)
		}

		if _, err := w.Write(j); err != nil {
			panic(err)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func solve(task goimp.Task) goimp.Result {
	var (
		freqs   = task.Freqs
		impData = task.ImpData
	)

	s := goimpcore.NewSolver(task.Code, freqs, impData)
	s.InitValues = task.InitValues
	s.SmartMode = "eis"

	result := s.Solve(10, 1000)

	return result
}
