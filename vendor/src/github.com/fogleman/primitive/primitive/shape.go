package primitive

import "github.com/fogleman/gg"

type Shape interface {
	Rasterize() []Scanline
	Copy() Shape
	Mutate()
	Draw(dc *gg.Context)
	SVG(attrs string) string
}
