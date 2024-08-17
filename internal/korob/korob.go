package korob

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path"
)

var (
	ErrOutputNotPNG = errors.New("png output required")
)

type Square struct {
	Row      int    `yaml:"row"`
	Column   int    `yaml:"column"`
	ColorHex string `yaml:"colorHex"`
}

type Korob struct {
	HorSquares     int      `yaml:"horSquares"`
	VerSquares     int      `yaml:"verSquares"`
	SquarePixels   int      `yaml:"squarePixels"`
	SquareColorHex string   `yaml:"squareColorHex"`
	MarginPixels   int      `yaml:"marginPixels"`
	MarginColorHex string   `yaml:"marginColorHex"`
	Grid           []Square `yaml:"grid"`

	pixWidth  int
	pixHeight int
	image     *image.RGBA
}

func hexToColor(hex string) color.RGBA {
	c := color.RGBA{A: 0xff}                                                //nolint:mnd // no alpha
	if n, _ := fmt.Sscanf(hex, "#%02x%02x%02x", &c.R, &c.G, &c.B); n != 3 { //nolint:mnd // RGB
		return color.RGBA{}
	}

	return c
}

func (k *Korob) Draw(output string) error {
	if path.Ext(output) != ".png" {
		return fmt.Errorf("%w: %s", ErrOutputNotPNG, output)
	}
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	k.pixWidth = k.HorSquares*k.SquarePixels + (k.HorSquares-1)*k.MarginPixels
	k.pixHeight = k.VerSquares*k.SquarePixels + (k.VerSquares-1)*k.MarginPixels
	k.image = image.NewRGBA(
		image.Rect(0, 0, k.pixWidth, k.pixHeight),
	)
	c := hexToColor(k.SquareColorHex)
	for i := range k.HorSquares {
		for j := range k.VerSquares {
			k.drawSquare(i, j, c)
		}
	}
	for _, s := range k.Grid {
		gc := hexToColor(s.ColorHex)
		k.drawSquare(s.Row, s.Column, gc)
	}
	mc := hexToColor(k.MarginColorHex)
	for i := range k.HorSquares - 1 {
		k.drawHorLine(i, mc)
	}
	for j := range k.VerSquares - 1 {
		k.drawVerLine(j, mc)
	}

	return png.Encode(file, k.image)
}

func (k *Korob) drawSquare(i, j int, paint color.RGBA) {
	pil := k.SquarePixels*i + i*k.MarginPixels
	pjl := k.SquarePixels*j + j*k.MarginPixels
	for x := pil; x < pil+k.SquarePixels; x++ {
		for y := pjl; y < pjl+k.SquarePixels; y++ {
			k.image.Set(x, y, paint)
		}
	}
}

func (k *Korob) drawHorLine(i int, paint color.RGBA) {
	pil := k.SquarePixels*(i+1) + k.MarginPixels*i
	for x := pil; x < pil+k.MarginPixels; x++ {
		for y := range k.pixHeight {
			k.image.Set(x, y, paint)
		}
	}
}

func (k *Korob) drawVerLine(j int, paint color.RGBA) {
	pjl := k.SquarePixels*(j+1) + k.MarginPixels*j
	for x := range k.pixWidth {
		for y := pjl; y < pjl+k.MarginPixels; y++ {
			k.image.Set(x, y, paint)
		}
	}
}
