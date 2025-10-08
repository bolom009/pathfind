package main

import (
	"context"
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/bolom009/delaunay"
	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind"
	"github.com/bolom009/pathfind/demo/utils"
	"github.com/bolom009/pathfind/graphs"
	"github.com/bolom009/pathfind/graphs/recast"
	"github.com/bolom009/pathfind/mesh"
	rlgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const floorPlan = `{"polygons":[{"outer":[{"x":276.63,"y":-251.83},{"x":309.33,"y":-253.34},{"x":335.9,"y":-178.9},{"x":388.6,"y":-146.17},{"x":454.03,"y":-146.17},{"x":454.03,"y":-111.73},{"x":486.73,"y":-113.24},{"x":513.74,"y":-37.58},{"x":481.04,"y":-36.06},{"x":481.04,"y":1.9},{"x":341.29,"y":1.9},{"x":341.29,"y":-5.98},{"x":291.13,"y":1.21},{"x":297.64,"y":19.42},{"x":264.94,"y":20.94},{"x":264.94,"y":58.9},{"x":125.19,"y":58.9},{"x":125.19,"y":44.37},{"x":85.44,"y":44.37},{"x":90.19,"y":57.68},{"x":57.5,"y":59.19},{"x":57.5,"y":97.16},{"x":-82.25,"y":97.16},{"x":-82.25,"y":23.45},{"x":-54.07,"y":23.45},{"x":-54.07,"y":59.19},{"x":-27.06,"y":-16.48},{"x":-54.65,"y":-16.48},{"x":-54.65,"y":-50.91},{"x":-18.27,"y":-50.91},{"x":-18.72,"y":-69.19},{"x":-28.3,"y":-69.19},{"x":-45.8,"y":-99.5},{"x":-28.3,"y":-129.81},{"x":6.7,"y":-129.81},{"x":8.84,"y":-126.1},{"x":163.89,"y":-180.78},{"x":163.89,"y":-211.9},{"x":192.07,"y":-211.9},{"x":192.07,"y":-176.16},{"x":219.08,"y":-251.83},{"x":191.49,"y":-251.83},{"x":191.49,"y":-286.27},{"x":276.63,"y":-286.27}],"innerHoles":[{"points":[{"x":23.45,"y":-100.79},{"x":24.2,"y":-99.5},{"x":9.88,"y":-74.7},{"x":10.46,"y":-50.91},{"x":30.49,"y":-50.91},{"x":30.49,"y":-16.48},{"x":63.19,"y":-17.99},{"x":75.19,"y":15.63},{"x":125.19,"y":15.63},{"x":125.19,"y":-14.8},{"x":153.37,"y":-14.8},{"x":153.37,"y":20.94},{"x":180.38,"y":-54.73},{"x":152.79,"y":-54.73},{"x":152.79,"y":-89.17},{"x":184.98,"y":-89.17},{"x":201.57,"y":-138.2},{"x":163.89,"y":-138.2},{"x":163.89,"y":-150.32}],"viewable":true},{"points":[{"x":303.64,"y":-138.2},{"x":231.92,"y":-138.2},{"x":215.32,"y":-89.17},{"x":237.93,"y":-89.17},{"x":237.93,"y":-54.73},{"x":270.63,"y":-56.24},{"x":281.27,"y":-26.41},{"x":341.29,"y":-35.02},{"x":341.29,"y":-71.8},{"x":369.47,"y":-71.8},{"x":369.47,"y":-36.06},{"x":396.48,"y":-111.73},{"x":368.89,"y":-111.73},{"x":368.89,"y":-124.58},{"x":303.64,"y":-165.11}],"viewable":true},{"points":[{"x":-25.71,"y":-99.5},{"x":-18.26,"y":-86.59},{"x":-3.34,"y":-86.59},{"x":4.11,"y":-99.5},{"x":-3.34,"y":-112.41},{"x":-18.26,"y":-112.41}],"viewable":true}],"obstacles":[{"points":[{"x":33.92,"y":54.61},{"x":23.68,"y":54.61},{"x":23.68,"y":15.99},{"x":33.92,"y":15.99}],"viewable":true},{"points":[{"x":383.1,"y":-41.91},{"x":338.2,"y":-41.46},{"x":338.1,"y":-51.69},{"x":383.0,"y":-52.14}],"viewable":true},{"points":[{"x":-5.91,"y":20.15},{"x":-27.36,"y":52.27},{"x":-30.09,"y":50.45},{"x":-8.64,"y":18.33}],"viewable":true},{"points":[{"x":25.56,"y":-47.41},{"x":27.21,"y":-45.76},{"x":27.81,"y":-43.5},{"x":27.21,"y":-41.24},{"x":25.56,"y":-39.59},{"x":23.3,"y":-38.99},{"x":21.04,"y":-39.59},{"x":19.39,"y":-41.24},{"x":18.79,"y":-43.5},{"x":19.39,"y":-45.76},{"x":21.04,"y":-47.41},{"x":23.3,"y":-48.01}],"viewable":false},{"points":[{"x":-40.05,"y":-25.45},{"x":-39.01,"y":-24.42},{"x":-38.63,"y":-23.0},{"x":-39.01,"y":-21.58},{"x":-40.05,"y":-20.55},{"x":-41.47,"y":-20.17},{"x":-47.13,"y":-20.17},{"x":-48.55,"y":-20.55},{"x":-49.59,"y":-21.58},{"x":-49.97,"y":-23.0},{"x":-49.59,"y":-24.42},{"x":-48.55,"y":-25.45},{"x":-47.13,"y":-25.83},{"x":-41.47,"y":-25.83}],"viewable":false},{"points":[{"x":-53.71,"y":63.09},{"x":-46.69,"y":70.11},{"x":-44.12,"y":79.7},{"x":-46.69,"y":89.29},{"x":-53.71,"y":96.31},{"x":-63.3,"y":98.88},{"x":-72.89,"y":96.31},{"x":-79.91,"y":89.29},{"x":-82.48,"y":79.7},{"x":-79.91,"y":70.11},{"x":-72.89,"y":63.09},{"x":-63.3,"y":60.52}],"viewable":false},{"points":[{"x":-16.91,"y":-0.94},{"x":-15.16,"y":0.81},{"x":-14.52,"y":3.2},{"x":-15.16,"y":5.59},{"x":-16.91,"y":7.34},{"x":-19.3,"y":7.98},{"x":-21.69,"y":7.34},{"x":-23.44,"y":5.59},{"x":-24.08,"y":3.2},{"x":-23.44,"y":0.81},{"x":-21.69,"y":-0.94},{"x":-19.3,"y":-1.58}],"viewable":false},{"points":[{"x":9.57,"y":3.33},{"x":11.39,"y":7.7},{"x":9.57,"y":12.07},{"x":5.2,"y":13.89},{"x":0.83,"y":12.07},{"x":-0.99,"y":7.7},{"x":0.83,"y":3.33},{"x":5.2,"y":1.51}],"viewable":true},{"points":[{"x":3.42,"y":58.28},{"x":5.05,"y":58.94},{"x":13.36,"y":64.88},{"x":11.62,"y":70.57},{"x":5.18,"y":78.27},{"x":1.63,"y":78.85},{"x":-5.03,"y":79.37},{"x":-13.36,"y":73.9},{"x":-12.08,"y":64.58},{"x":-11.45,"y":63.54},{"x":-6.37,"y":58.62},{"x":-3.54,"y":58.03}],"viewable":false},{"points":[{"x":204.92,"y":-4.72},{"x":206.55,"y":-4.06},{"x":214.86,"y":1.88},{"x":213.12,"y":7.57},{"x":206.68,"y":15.27},{"x":203.13,"y":15.85},{"x":196.47,"y":16.37},{"x":188.14,"y":10.9},{"x":189.42,"y":1.58},{"x":190.05,"y":0.54},{"x":195.13,"y":-4.38},{"x":197.96,"y":-4.97}],"viewable":false},{"points":[{"x":51.84,"y":5.47},{"x":49.4,"y":7.67},{"x":44.4,"y":7.67},{"x":41.96,"y":5.47},{"x":42.15,"y":-5.5},{"x":51.65,"y":-5.5}],"viewable":false}]},{"outer":[{"x":-90.95,"y":-304.37},{"x":-86.3,"y":-305.79},{"x":-77.93,"y":-278.3},{"x":-82.34,"y":-276.95},{"x":-47.07,"y":-164.58},{"x":-45.02,"y":-165.31},{"x":-35.46,"y":-138.22},{"x":-205.38,"y":-78.29},{"x":-214.94,"y":-105.38},{"x":-211.7,"y":-106.51},{"x":-249.25,"y":-226.12},{"x":-250.3,"y":-225.81},{"x":-258.67,"y":-253.3},{"x":-230.47,"y":-262.26},{"x":-230.36,"y":-261.91},{"x":-118.44,"y":-296.0},{"x":-120.19,"y":-301.55},{"x":-92.77,"y":-310.16}],"innerHoles":[{"points":[{"x":-221.75,"y":-234.5},{"x":-184.58,"y":-116.08},{"x":-74.19,"y":-155.01},{"x":-109.84,"y":-268.58}],"viewable":true}],"obstacles":[]},{"outer":[{"x":32.0,"y":-200.2},{"x":14.5,"y":-169.89},{"x":-20.5,"y":-169.89},{"x":-38.0,"y":-200.2},{"x":-20.5,"y":-230.51},{"x":14.5,"y":-230.51}],"innerHoles":[{"points":[{"x":-17.91,"y":-200.2},{"x":-10.46,"y":-187.29},{"x":4.46,"y":-187.29},{"x":11.91,"y":-200.2},{"x":4.46,"y":-213.11},{"x":-10.46,"y":-213.11}],"viewable":true}],"obstacles":[]},{"outer":[{"x":129.6,"y":-227.3},{"x":112.1,"y":-196.99},{"x":77.1,"y":-196.99},{"x":59.6,"y":-227.3},{"x":77.1,"y":-257.61},{"x":112.1,"y":-257.61}],"innerHoles":[{"points":[{"x":79.69,"y":-227.3},{"x":87.14,"y":-214.39},{"x":102.06,"y":-214.39},{"x":109.51,"y":-227.3},{"x":102.06,"y":-240.21},{"x":87.14,"y":-240.21}],"viewable":true}],"obstacles":[{"points":[{"x":97.06,"y":-209.01},{"x":98.71,"y":-207.36},{"x":99.31,"y":-205.1},{"x":98.71,"y":-202.84},{"x":97.06,"y":-201.19},{"x":94.8,"y":-200.59},{"x":92.54,"y":-201.19},{"x":90.89,"y":-202.84},{"x":90.29,"y":-205.1},{"x":90.89,"y":-207.36},{"x":92.54,"y":-209.01},{"x":94.8,"y":-209.61}],"viewable":false}]}]}`

func main() {
	location, err := utils.NewLocationDataFromJSON([]byte(floorPlan))
	if err != nil {
		panic(err)
	}

	rPolygons := make([]*mesh.Polygon, len(location.Polygons))
	for j, polygon := range location.Polygons {
		innerHoles := make([]*mesh.Hole, len(polygon.InnerHoles))
		for i, innerHole := range polygon.InnerHoles {
			innerHoles[i] = mesh.NewInnerHole(innerHole.Points, 3)
		}

		obstacles := make([]*mesh.Hole, len(polygon.Obstacles))
		for i, obstacle := range polygon.Obstacles {
			obstacles[i] = mesh.NewObstacle(obstacle.Points, 3, obstacle.Viewable)
		}

		rPolygons[j] = mesh.NewPolygon(polygon.Outer, innerHoles, obstacles, 3)
	}

	var (
		screen      = geom.Vector2{800, 600}
		camera      = rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(-screen.X/2, -screen.Y/2), 0, 0.5)
		path        = make([]geom.Vector2, 0)
		start       = geom.Vector2{X: 202, Y: -268}
		dest        = geom.Vector2{X: -4, Y: 84}
		recastGraph = recast.NewRecast(rPolygons, recast.WithSearchOutOfArea(true))
		pathfinder  = pathfind.NewPathfinder[geom.Vector2]([]graphs.NavGraph[geom.Vector2]{
			recastGraph,
		})
		graphId        = 0
		isDrawGraph    = true
		edgesCount     = 0
		vertexCount    = 0
		trianglesCount = 0

		initTime string
		pathTime string
	)

	t := time.Now()
	if err := pathfinder.Initialize(context.Background()); err != nil {
		panic(err)
	}
	initTime = time.Since(t).String()

	t = time.Now()
	path = pathfinder.Path(graphId, start, dest)
	pathTime = time.Since(t).String()

	visGraph := pathfinder.Graph(graphId)
	for _, edges := range visGraph {
		vertexCount++
		edgesCount += len(edges)
	}

	trianglesCount = len(recastGraph.Triangles())
	triangles := recastGraph.Triangles()

	bgColor := rl.NewColor(176, 165, 158, 255)
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
		rl.ClearBackground(bgColor)

		if !rl.IsKeyDown(rl.KeyLeftControl) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			dest = geom.Vector2{X: mouseWorldPos.X, Y: -mouseWorldPos.Y}
			t := time.Now()
			searchPath := pathfinder.Path(graphId, start, dest)
			path = searchPath
			pathTime = time.Since(t).String()
		}

		if rl.IsKeyDown(rl.KeyLeftControl) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			start = geom.Vector2{X: mouseWorldPos.X, Y: -mouseWorldPos.Y}
			t := time.Now()
			searchPath := pathfinder.Path(graphId, start, dest)
			path = searchPath
			pathTime = time.Since(t).String()
		}

		// drawing map
		for _, polygon := range location.Polygons {
			drawMap(polygon.Outer, polygon.InnerHoles, rl.White, rl.Blue)
			for _, obstacle := range polygon.Obstacles {
				drawArea(obstacle.Points, rl.Gray)
			}
		}

		drawTriangles(triangles, rl.LightGray)

		//if isDrawGraph {
		//	drawGraph(pathfinder.GraphWithSearchPath(graphId, start, dest))
		//}
		drawPath(start, dest, path, camera.Zoom, true)

		rl.EndMode2D()

		drawTopPanel(int32(screen.X), rl.Vector2{}, &isDrawGraph, initTime, pathTime, trianglesCount, vertexCount, edgesCount)
		rl.SetWindowTitle(fmt.Sprintf("Test A* (%v, %v)", int(mouseWorldPos.X), int(-mouseWorldPos.Y)))

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

	rl.DrawCircleV(rl.NewVector2(start.X, -start.Y), 3/zoom, color.RGBA{R: 0x90, G: 0xee, B: 0x90, A: 0xff})
	rl.DrawCircleV(rl.NewVector2(dest.X, -dest.Y), 3/zoom, color.RGBA{R: 0xe7, G: 0x6f, B: 0x51, A: 0xFF})

	if len(path) == 0 {
		return
	}

	for i := range len(path) - 1 {
		p1, p2 := path[i], path[i+1]
		rl.DrawLineV(rl.NewVector2(p1.X, -p1.Y), rl.NewVector2(p2.X, -p2.Y), color.RGBA{R: 0x2a, G: 0x9d, B: 0x8f, A: 0xFF})
		if !isSkipNumbers {
			p := p1.Add(p2).Div(2)
			rl.DrawText(fmt.Sprintf("%v", i+1), int32(p.X), int32(-p.Y), 10, rl.Red)
		}
	}
}

func drawTriangles(triangles []recast.Triangle, color rl.Color) {
	for _, pp := range triangles {
		rl.DrawTriangleFan([]rl.Vector2{
			rl.NewVector2(pp[0].X, -pp[0].Y),
			rl.NewVector2(pp[1].X, -pp[1].Y),
			rl.NewVector2(pp[2].X, -pp[2].Y),
		}, color)
	}
}

func drawMap(polygon []geom.Vector2, holes []utils.Hole, color rl.Color, holesColor rl.Color) {
	drawArea(polygon, color)
	for _, hole := range holes {
		drawArea(hole.Points, holesColor)
	}
}

func drawArea(polygon []geom.Vector2, color rl.Color) {
	pp := make([]*delaunay.Point, 0)
	for _, p := range polygon {
		pp = append(pp, delaunay.NewPoint(p.X, p.Y))
	}

	d := delaunay.NewSweepContext(pp, nil)
	dTriangles := d.Triangulate()
	triangles := make([]recast.Triangle, 0)
	for _, t := range dTriangles {
		dpp := t.GetPoints()
		triangles = append(triangles, recast.Triangle{
			geom.Vector2{dpp[0].X, dpp[0].Y},
			geom.Vector2{dpp[1].X, dpp[1].Y},
			geom.Vector2{dpp[2].X, dpp[2].Y},
		})
	}

	drawTriangles(triangles, color)
}

func drawGraph(graph map[geom.Vector2][]geom.Vector2) {
	for p, elems := range graph {
		for _, elem := range elems {
			rl.DrawLine(int32(p.X), int32(p.Y), int32(elem.X), int32(elem.Y), rl.NewColor(230, 41, 55, 30))
		}
	}
}

func drawTopPanel(width int32, tPos rl.Vector2, isDrawGraph *bool, initTime,
	pathTime string, trianglesCount, vertexCount, edgesCount int) {

	rl.DrawRectangle(int32(tPos.X), int32(tPos.Y), width, 30, rl.NewColor(255, 255, 255, 100))

	*isDrawGraph = rlgui.CheckBox(rl.NewRectangle(15, 10, 15, 15), "draw graph", *isDrawGraph)

	rl.DrawText(" | ", 200, 10, 15, rl.Gray)
	rlgui.Label(rl.NewRectangle(220, 10, 150, 15), "Init time: "+initTime)
	rl.DrawText(" | ", 300, 10, 15, rl.Gray)
	rlgui.Label(rl.NewRectangle(320, 10, 150, 15), "Path time: "+pathTime)
	rl.DrawText(" | ", 420, 10, 15, rl.Gray)
	rlgui.Label(rl.NewRectangle(440, 10, 200, 15), fmt.Sprintf("Triangles %v", trianglesCount))
	rl.DrawText(" | ", 520, 10, 15, rl.Gray)
	rlgui.Label(rl.NewRectangle(540, 10, 180, 15), fmt.Sprintf("Graph (%vx%v)", vertexCount, edgesCount))
}
