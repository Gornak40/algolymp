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

	image *image.RGBA
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

	k.image = image.NewRGBA(
		image.Rect(
			0, 0,
			k.HorSquares*k.SquarePixels+(k.HorSquares-1)*k.MarginPixels,
			k.VerSquares*k.SquarePixels+(k.VerSquares-1)*k.MarginPixels,
		),
	)
	for i := range k.HorSquares {
		for j := range k.VerSquares {
			k.drawSquare(i, j, color.RGBA{255, 0, 0, 255})
		}
	}
	for _, s := range k.Grid {
		k.drawSquare(s.Row, s.Column, color.RGBA{0, 255, 0, 255})
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
