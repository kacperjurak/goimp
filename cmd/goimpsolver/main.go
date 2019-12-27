package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/kacperjurak/goimp/cmd"
	"github.com/kacperjurak/goimpcore"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		code       string
		file       string
		initValues cmd.ArrayFlags
		cutLow     uint
		cutHigh    uint
		unity      bool
		flip       bool
		imgOut     bool
		imgSave    bool
		imgPath    string
		imgDPI     uint
		imgSize    uint
		quiet      bool
		flipY      = -1.0
	)

	flag.StringVar(&code, "c", "", "Boukamp Circuit Description Code")
	flag.StringVar(&file, "f", "data.txt", "Measurement data file")
	flag.Var(&initValues, "v", "Parameters init values (array)")
	flag.UintVar(&cutLow, "b", 0, "Cut X of begining frequencies from a file")
	flag.UintVar(&cutHigh, "e", 0, "Cut X of ending frequencies from a file")
	flag.BoolVar(&unity, "unity", false, "Use Unity weighting intead Modulus")
	flag.BoolVar(&flip, "noflip", false, "Don't flip imaginary part on image")
	flag.BoolVar(&imgOut, "imgout", false, "Image data to STDOUT")
	flag.BoolVar(&imgSave, "imgsave", false, "Save image to file")
	flag.StringVar(&imgPath, "imgpath", "eis.svg", "Path to generated image")
	flag.UintVar(&imgDPI, "dpi", 96, "Image DPI")
	flag.UintVar(&imgSize, "imgsize", 4, "Image size (inches)")
	flag.BoolVar(&quiet, "q", false, "Quiet mode")
	flag.Parse()

	freqs, impData := parseFile(file)

	freqs = freqs[cutLow : len(freqs)-int(cutHigh)]
	impData = impData[cutLow : len(impData)-int(cutHigh)]

	s := goimpcore.NewSolver(code, initValues, goimpcore.MODULUS)

	res, err := s.Solve(freqs, impData)
	if err != nil {
		panic(err)
	}

	if imgOut || imgSave {
		l, err := vg.ParseLength(strconv.Itoa(int(imgSize*imgDPI)) + "pt")
		if err != nil {
			panic(err)
		}

		p, err := plot.New()
		if err != nil {
			panic(err)
		}

		p.Title.Text = "ChiSq: " + fmt.Sprintf("%e", res.ChiSq)
		p.X.Label.Text = "Zr"
		p.Y.Label.Text = "Zi"

		ptsCalc := make(plotter.XYs, len(impData))
		ptsMeas := make(plotter.XYs, len(impData))

		if flip {
			flipY = 1.0
		}

		calculated := goimpcore.CircuitImpedance(code, freqs, res.Params)
		for i, p := range impData {
			ptsCalc[i].X = calculated[i][0]
			ptsCalc[i].Y = calculated[i][1] * flipY
			ptsMeas[i].X = p[0]
			ptsMeas[i].Y = p[1] * flipY
		}

		err = plotutil.AddLinePoints(p, "calculated", ptsCalc, "measured", ptsMeas)
		if err != nil {
			panic(err)
		}

		if imgSave {
			if err := p.Save(l, l, imgPath); err != nil {
				panic(err)
			}
		}

		if imgOut {
			c := vgsvg.New(l, l)
			p.Draw(draw.New(c))
			if _, err := c.WriteTo(os.Stdout); err != nil {
				panic(err)
			}
		}
	}
	if !(quiet) {
		log.Println("Result:", res)
	}
}

func parseFile(file string) (freqs []float64, impData [][2]float64) {
	f, _ := os.Open(file)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		var lineVals [3]float64

		for i := 0; i < 3; i++ {
			l := strings.Fields(line)
			val, err := strconv.ParseFloat(l[i], 64)
			if err != nil {
				log.Fatal(err)
			}
			lineVals[i] = val
		}
		//measData = append(measData, lineVals)
		freqs = append(freqs, lineVals[0])
		impData = append(impData, [2]float64{lineVals[1], lineVals[2]})
	}
	return freqs, impData
}
