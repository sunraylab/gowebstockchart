package stockchart

import (
	"github.com/gowebapi/webapi/core/js"
	"github.com/gowebapi/webapi/html/canvas"
	"github.com/larry868/datarange"
	"github.com/larry868/rgb"
)

type DrawingYGrid struct {
	Drawing
	fScale     bool // Draw the scale, otherwise only the lines
	lastyrange datarange.DataRange
}

func NewDrawingYGrid(series *DataList, fscale bool) *DrawingYGrid {
	drawing := new(DrawingYGrid)
	drawing.Name = "ygrid"
	drawing.series = series
	drawing.MainColor = rgb.Gray.Lighten(0.85)
	drawing.fScale = fscale

	drawing.Drawing.OnRedraw = func() {
		//	yrange := drawing.series.DataRange(drawing.xAxisRange, 10)
		drawing.chart.yAxisRange = drawing.series.DataRange(&drawing.chart.selectedTimeSlice, 10)
		drawing.lastyrange = drawing.chart.yAxisRange
		// Debug(DBG_REDRAW, "%q OnRedraw drawarea:%s, xAxisRange:%v, datarange:%v", drawing.Name, drawing.drawArea, drawing.xAxisRange.String(), drawing.chart.yAxisRange)
		drawing.onRedraw()
	}
	drawing.Drawing.NeedRedraw = func() bool {
		ynewrange := drawing.series.DataRange(&drawing.chart.selectedTimeSlice, 10)
		return !ynewrange.Equal(drawing.lastyrange)
	}
	return drawing
}

// OnRedraw redraw the Y axis
func (drawing DrawingYGrid) onRedraw() {

	// setup default text drawing properties
	drawing.Ctx2D.SetTextAlign(canvas.StartCanvasTextAlign)
	drawing.Ctx2D.SetTextBaseline(canvas.MiddleCanvasTextBaseline)
	drawing.Ctx2D.SetFont(`12px 'Roboto', sans-serif`)

	// draw the Y Scale
	yrange := drawing.chart.yAxisRange
	for val := yrange.High(); val >= yrange.Low() && yrange.StepSize() > 0; val -= yrange.StepSize() {

		// calculate ypos
		yrate := yrange.Progress(val)
		ypos := float64(drawing.drawArea.End().Y) - yrate*float64(drawing.drawArea.Height)
		ypos = float64(drawing.drawArea.BoundY(int(ypos)))

		// draw the grid line
		drawing.Ctx2D.SetFillStyle(&canvas.Union{Value: js.ValueOf(drawing.MainColor.Hexa())})
		linew := 10.0
		if !drawing.fScale {
			linew = float64(drawing.drawArea.Width)
		}
		drawing.Ctx2D.FillRect(float64(drawing.drawArea.O.X), ypos, linew, 1)

		// draw yscale label
		if drawing.fScale {
			strvalue := datarange.FormatData(val, yrange.StepSize()) // fmt.Sprintf("%.1f", val)
			drawing.Ctx2D.SetFillStyle(&canvas.Union{Value: js.ValueOf(rgb.Gray.Darken(0.5).Hexa())})
			drawing.Ctx2D.FillText(strvalue, float64(drawing.drawArea.O.Y+7), ypos+1, nil)
		}
	}
}
