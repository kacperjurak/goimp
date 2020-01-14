package main

import (
	"flag"
	"fmt"
	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/kacperjurak/goimp/cmd"
	"github.com/kacperjurak/goimpcore"
	"github.com/kacperjurak/gologspace"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"gonum.org/v1/plot/vg/vgsvg"
	"image"
	"image/color"
	"os"
	"strconv"
)

var (
	code        string
	values      cmd.ArrayFlags
	freqs       cmd.ArrayFlags
	flip        bool
	imgOut      bool
	imgTerm     bool
	imgTermSize uint
	imgSave     bool
	imgPath     string
	imgDPI      uint
	imgSize     uint
	freqAuto    bool
	freqMin     float64
	freqMax     float64
	pointsNo    uint
	quiet       bool
	littleNoise bool
	noisyPoints uint
	noiseLevel  float64
	res         [][2]float64
)

func main() {
	flag.StringVar(&code, "c", "", "Boukamp Circuit Description Code")
	flag.Var(&values, "v", "Parameter values (array)")
	flag.Var(&freqs, "f", "Frequencies (Hz) (array)")
	flag.BoolVar(&flip, "noflip", false, "Don't flip imaginary part on image")
	flag.BoolVar(&imgOut, "imgout", false, "Image data to STDOUT")
	flag.BoolVar(&imgTerm, "imgterm", false, "Show image on terminal")
	flag.UintVar(&imgTermSize, "termsize", 320, "Terminal size")
	flag.BoolVar(&imgSave, "imgsave", false, "Save image to file?")
	flag.StringVar(&imgPath, "imgpath", "eis.svg", "Path to generated image")
	flag.UintVar(&imgDPI, "dpi", 96, "Image DPI")
	flag.UintVar(&imgSize, "imgsize", 4, "Image size (inches)")
	flag.BoolVar(&freqAuto, "fauto", true, "Auto frequencies")
	flag.Float64Var(&freqMin, "fmin", 0.001, "Min frequency")
	flag.Float64Var(&freqMax, "fmax", 1000000, "Max frequency")
	flag.UintVar(&pointsNo, "fpn", 100, "Number of frequency points")
	flag.BoolVar(&quiet, "q", false, "Quiet mode")
	flag.BoolVar(&littleNoise, "n", false, "Add a little noise")
	flag.UintVar(&noisyPoints, "np", 0, "Number of noisy points")
	flag.Float64Var(&noiseLevel, "nl", 0.3, "Noise level (1 = 100%)")
	flag.Parse()

	if freqAuto {
		freqs = gologspace.Generate(freqMin, freqMax, pointsNo)
	}

	res = goimpcore.CircuitImpedanceNoisy(code, freqs, values, noisyPoints, noiseLevel, littleNoise)

	if !flip {
		for i, v := range res {
			res[i] = [2]float64{v[0], v[1] * -1}
		}
	}

	if imgOut || imgSave || imgTerm {
		l, err := vg.ParseLength(strconv.Itoa(int(imgSize*imgDPI)) + "pt")
		if err != nil {
			panic(err)
		}

		p, err := plot.New()
		if err != nil {
			panic(err)
		}

		p.X.Label.Text = "Zr"
		p.Y.Label.Text = "Zi"

		err = plotutil.AddLinePoints(p, "", points(res))
		if err != nil {
			panic(err)
		}

		if imgSave {
			if err := p.Save(l, l, imgPath); err != nil {
				panic(err)
			}
		}

		if imgTerm {
			img := image.NewRGBA(image.Rect(0, 0, int(imgSize*imgDPI), int(imgSize*imgDPI)))
			c := vgimg.NewWith(vgimg.UseImage(img))
			p.Draw(draw.New(c))
			pix, err := ansimage.NewScaledFromImage(img, int(imgTermSize), int(imgTermSize), color.Opaque, ansimage.ScaleModeFit, ansimage.DitheringWithBlocks)
			if err != nil {
				panic(err)
			}
			ansimage.ClearTerminal()
			pix.DrawExt(false, true)
		}

		if imgOut {
			c := vgsvg.New(l, l)
			p.Draw(draw.New(c))
			if _, err := c.WriteTo(os.Stdout); err != nil {
				panic(err)
			}
		}
	}

	if !(imgOut || imgTerm || quiet) {
		fmt.Print(VSlice{
			freqs: freqs,
			res:   res},
		)
	}
}

func points(data [][2]float64) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X = data[i][0]
		pts[i].Y = data[i][1]
	}
	return pts
}

type VSlice struct {
	freqs []float64
	res   [][2]float64
}

func (s VSlice) String() string {
	var r string
	for i, v := range s.res {
		r += fmt.Sprintln(s.freqs[i], v[0], v[1])
	}
	return r
}
