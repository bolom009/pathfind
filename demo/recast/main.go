package main

import (
	"context"
	"fmt"
	"github.com/bolom009/clipper"
	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind"
	"github.com/bolom009/pathfind/demo/recast/polyjson"
	"github.com/bolom009/pathfind/graphs"
	"github.com/bolom009/pathfind/graphs/recast"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
	"math"
	"time"
)

// Created with Polygon Constructor: https://alaricus.github.io/PolygonConstructor/
// const floorPlan = `{"canvas":{"w":800,"h":600},"polygons":[[{"x":15.4375,"y":582},{"x":13.4375,"y":15},{"x":784.4375,"y":18},{"x":790.4375,"y":587},{"x":750.4375,"y":588},{"x":740.4375,"y":61},{"x":656.4375,"y":61},{"x":655.4375,"y":117},{"x":694.4375,"y":117},{"x":695.4375,"y":154},{"x":655.4375,"y":152},{"x":658.4375,"y":211},{"x":701.4375,"y":209},{"x":702.4375,"y":263},{"x":661.4375,"y":262},{"x":663.4375,"y":309},{"x":705.4375,"y":309},{"x":708.4375,"y":357},{"x":667.4375,"y":358},{"x":670.4375,"y":406},{"x":714.4375,"y":406},{"x":715.4375,"y":463},{"x":610.4375,"y":464},{"x":609.4375,"y":406},{"x":641.4375,"y":407},{"x":637.4375,"y":357},{"x":606.4375,"y":358},{"x":605.4375,"y":311},{"x":638.4375,"y":309},{"x":634.4375,"y":261},{"x":602.4375,"y":261},{"x":599.4375,"y":207},{"x":638.4375,"y":208},{"x":629.4375,"y":151},{"x":592.4375,"y":151},{"x":594.4375,"y":115},{"x":636.4375,"y":114},{"x":628.4375,"y":61},{"x":583.4375,"y":66},{"x":517.4375,"y":65},{"x":521.4375,"y":466},{"x":582.4375,"y":465},{"x":583.4375,"y":512},{"x":718.4375,"y":513},{"x":722.4375,"y":583},{"x":525.4375,"y":582},{"x":450.4375,"y":584},{"x":448.4375,"y":466},{"x":498.4375,"y":469},{"x":491.4375,"y":64},{"x":434.4375,"y":64},{"x":448.4375,"y":427},{"x":473.4375,"y":428},{"x":463.4375,"y":90},{"x":480.4375,"y":88},{"x":484.4375,"y":446},{"x":412.4375,"y":447},{"x":366.4375,"y":580},{"x":202.4375,"y":582},{"x":48.4375,"y":582}],[{"x":141.4375,"y":68},{"x":144.4375,"y":189},{"x":47.4375,"y":186},{"x":44.4375,"y":69}],[{"x":368.4375,"y":65},{"x":370.4375,"y":190},{"x":265.4375,"y":184},{"x":262.4375,"y":68}],[{"x":374.4375,"y":268},{"x":378.4375,"y":374},{"x":274.4375,"y":374},{"x":268.4375,"y":279}],[{"x":153.4375,"y":281},{"x":159.4375,"y":380},{"x":55.4375,"y":377},{"x":51.4375,"y":280}],[{"x":241.4375,"y":207},{"x":242.4375,"y":259},{"x":191.4375,"y":261},{"x":188.4375,"y":209}],[{"x":216.4375,"y":109},{"x":219.4375,"y":135},{"x":191.4375,"y":136},{"x":191.4375,"y":113}],[{"x":232.4375,"y":327},{"x":233.4375,"y":358},{"x":198.4375,"y":359},{"x":196.4375,"y":331}],[{"x":51.4375,"y":432},{"x":372.4375,"y":429},{"x":314.4375,"y":543},{"x":25.4375,"y":546}],[{"x":387.4375,"y":470},{"x":353.4375,"y":551},{"x":335.4375,"y":541},{"x":378.4375,"y":460}],[{"x":376.4375,"y":392},{"x":370.4375,"y":418},{"x":278.4375,"y":417},{"x":281.4375,"y":401}]]}`
const floorPlan = `{"canvas":{"w":800,"h":600},"polygons":[[{"x":0,"y":0},{"x":120,"y":0},{"x":120,"y":340},{"x":180,"y":340},{"x":180,"y":-120},{"x":300,"y":-120},{"x":300,"y":340},{"x":360,"y":340},{"x":360,"y":0},{"x":480,"y":0},{"x":480,"y":420},{"x":300,"y":420},{"x":300,"y":540},{"x":340,"y":540},{"x":340,"y":720},{"x":140,"y":720},{"x":140,"y":540},{"x":180,"y":540},{"x":180,"y":420},{"x":0,"y":420}]]}`

func main() {
	polygon, holes, screen, err := polyjson.NewPolygonsFromJSON([]byte(floorPlan))
	if err != nil {
		panic(err)
	}

	nPolygon := clipper.OffsetPolygon(polygon, 20.0)
	nHoles := make([][]geom.Vector2, len(holes))
	for i, hole := range holes {
		nHoles[i] = clipper.OffsetPolygon(hole, -5.0)
	}

	var (
		camera      = rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(-screen.X/2, -screen.Y/2), 0, 0.5)
		path        = make([]geom.Vector2, 0)
		start       = geom.Vector2{X: 25, Y: 25}
		dest        = geom.Vector2{X: 430, Y: 200}
		recastGraph = recast.NewRecast(nPolygon, nHoles)
		pathfinder  = pathfind.NewPathfinder[geom.Vector2]([]graphs.NavGraph[geom.Vector2]{
			recastGraph,
		})
		graphId = 0
	)

	if err := pathfinder.Initialize(context.Background()); err != nil {
		panic(err)
	}

	path = pathfinder.Path(graphId, start, dest)

	rl.SetTraceLogLevel(rl.LogError)
	rl.SetTargetFPS(60)
	rl.InitWindow(int32(screen.X), int32(screen.Y), "Test *A")
	for {
		if rl.WindowShouldClose() {
			break
		}

		mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
		moveCamera(&camera, mouseWorldPos)

		rl.BeginDrawing()

		rl.BeginMode2D(camera)
		rl.ClearBackground(rl.White)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			dest = geom.Vector2{X: mouseWorldPos.X, Y: mouseWorldPos.Y}
			t := time.Now()
			searchPath := pathfinder.Path(graphId, start, dest)
			if len(searchPath) >= 2 {
				path = searchPath
			}
			fmt.Println("ttt", time.Since(t).String())
			i = 0
		}

		// drawing map
		drawMap(polygon, holes)
		drawGraph(pathfinder.GraphWithSearchPath(graphId, start, dest))
		drawPath(start, dest, path, camera.Zoom, true)

		rl.EndMode2D()

		rl.SetWindowTitle(fmt.Sprintf("Test A* (%v, %v)", int(mouseWorldPos.X), int(mouseWorldPos.Y)))

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func moveCamera(camera *rl.Camera2D, mouseWorldPos rl.Vector2) {
	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		// Set the offset to where the mouse is
		camera.Offset = rl.GetMousePosition()

		// Set the target to match, so that the camera maps the world space point
		// under the cursor to the screen space point under the cursor at any zoom
		camera.Target = mouseWorldPos

		// Zoom increment
		scaleFactor := 1.0 + (0.25 * math.Abs(float64(wheel)))
		if wheel < 0 {
			scaleFactor = 1.0 / scaleFactor
		}
		camera.Zoom = rl.Clamp(camera.Zoom*float32(scaleFactor), 0.125, 64.0)
	}

	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		delta := rl.GetMouseDelta()
		delta = rl.Vector2Scale(delta, -1.0/camera.Zoom)
		camera.Target = rl.Vector2Add(camera.Target, delta)
	}
}

func drawPath(start, dest geom.Vector2, path []geom.Vector2, zoom float32, skipNumbers ...bool) {
	isSkipNumbers := false
	if len(skipNumbers) > 0 {
		isSkipNumbers = true
	}

	rl.DrawCircle(int32(start.X), int32(start.Y), 3/zoom, color.RGBA{R: 0x90, G: 0xee, B: 0x90, A: 0xff})
	rl.DrawCircle(int32(dest.X), int32(dest.Y), 3/zoom, color.RGBA{R: 0xe7, G: 0x6f, B: 0x51, A: 0xFF})

	if len(path) == 0 {
		return
	}

	for i := range len(path) - 1 {
		p1, p2 := path[i], path[i+1]
		rl.DrawLine(int32(p1.X), int32(p1.Y), int32(p2.X), int32(p2.Y), color.RGBA{R: 0x2a, G: 0x9d, B: 0x8f, A: 0xFF})

		if !isSkipNumbers {
			p := p1.Add(p2).Div(2)
			rl.DrawText(fmt.Sprintf("%v", i+1), int32(p.X), int32(p.Y), 10, rl.Red)
		}
	}
}

func drawMap(polygon []geom.Vector2, holes [][]geom.Vector2) {
	polyLen := len(polygon)
	for range polygon {
		for i := range polyLen - 1 {
			p1, p2 := polygon[i], polygon[i+1]
			rl.DrawLine(int32(p1.X), int32(p1.Y), int32(p2.X), int32(p2.Y), rl.SkyBlue)
		}

		// DRAW END LINES
		startPos, endPos := polygon[polyLen-1], polygon[0]
		rl.DrawLine(int32(startPos.X), int32(startPos.Y), int32(endPos.X), int32(endPos.Y), rl.SkyBlue)
	}

	for _, hole := range holes {
		n := len(hole)
		if n < 2 {
			return
		}

		for i := range len(hole) - 1 {
			p1, p2 := hole[i], hole[i+1]
			rl.DrawLine(int32(p1.X), int32(p1.Y), int32(p2.X), int32(p2.Y), rl.SkyBlue)
		}

		// DRAW END LINES
		startPos, endPos := hole[n-1], hole[0]
		rl.DrawLine(int32(startPos.X), int32(startPos.Y), int32(endPos.X), int32(endPos.Y), rl.SkyBlue)
	}
}

var i = 0

func drawGraph(graph map[geom.Vector2][]geom.Vector2) {

	eCount := 0
	pCount := 0
	for p, elems := range graph {
		for _, elem := range elems {
			rl.DrawLine(int32(p.X), int32(p.Y), int32(elem.X), int32(elem.Y), rl.NewColor(230, 41, 55, 30))
		}
		eCount += len(elems)
		pCount++
	}

	//fmt.Println(pCount, eCount)
	i++
}
