package graph

import "github.com/fogleman/gg"

type LineStyle struct {
	solid   bool
	dots    bool
	pillars bool
}

func NewLS() *LineStyle {
	return &LineStyle{}
}

func (ls *LineStyle) IsSolid() bool {
	return ls.solid
}

func (ls *LineStyle) IsDots() bool {
	return ls.dots
}

func (ls *LineStyle) IsPillars() bool {
	return ls.pillars
}

func (ls *LineStyle) Solid() {
	ls.solid = true
}

func (ls *LineStyle) Dots() {
	ls.dots = true
}

func (ls *LineStyle) Pillars() {
	ls.pillars = true
}

func (ls *LineStyle) SetLineParams(dc *gg.Context) {
	if ls.solid {
		dc.SetLineWidth(2)
	}

	if ls.dots {
		dc.SetLineWidth(1)
	}

	if ls.pillars {
		dc.SetLineWidth(10)
		dc.SetLineCap(gg.LineCapSquare)
	}
}

func (ls *LineStyle) DrawLine(dc *gg.Context, x, y []float64, originY float64) {
	if ls.solid && len(x) > 0 {
		dc.NewSubPath()
		dc.MoveTo(x[0], y[0])
		for i := 1; i < len(x); i++ {
			dc.LineTo(x[i], y[i])
		}
		dc.Stroke()
	}

	if ls.dots {
		for i := range x {
			dc.DrawCircle(x[i], y[i], 5)
			dc.Fill()
		}
	}

	if ls.pillars {
		for i := range x {
			dc.NewSubPath()
			dc.MoveTo(x[i], originY)
			dc.LineTo(x[i], y[i])
			dc.Stroke()
		}
	}
}
