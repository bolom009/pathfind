package main

import (
	"context"
	"fmt"
	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind/demo/grid/polyjson"
	"image/color"
	"math"
	"time"

	"github.com/bolom009/pathfind"
	"github.com/bolom009/pathfind/graphs"
	"github.com/bolom009/pathfind/graphs/grid"
	"github.com/bolom009/pathfind/obstacles"
	rlgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	searchPathTimeInterval = time.Millisecond * 300
)

// Created with Polygon Constructor: https://alaricus.github.io/PolygonConstructor/
// const floorPlan = `{"canvas":{"w":800,"h":600},"polygons":[[{"x":15.4375,"y":582},{"x":13.4375,"y":15},{"x":784.4375,"y":18},{"x":790.4375,"y":587},{"x":750.4375,"y":588},{"x":740.4375,"y":61},{"x":656.4375,"y":61},{"x":655.4375,"y":117},{"x":694.4375,"y":117},{"x":695.4375,"y":154},{"x":655.4375,"y":152},{"x":658.4375,"y":211},{"x":701.4375,"y":209},{"x":702.4375,"y":263},{"x":661.4375,"y":262},{"x":663.4375,"y":309},{"x":705.4375,"y":309},{"x":708.4375,"y":357},{"x":667.4375,"y":358},{"x":670.4375,"y":406},{"x":714.4375,"y":406},{"x":715.4375,"y":463},{"x":610.4375,"y":464},{"x":609.4375,"y":406},{"x":641.4375,"y":407},{"x":637.4375,"y":357},{"x":606.4375,"y":358},{"x":605.4375,"y":311},{"x":638.4375,"y":309},{"x":634.4375,"y":261},{"x":602.4375,"y":261},{"x":599.4375,"y":207},{"x":638.4375,"y":208},{"x":629.4375,"y":151},{"x":592.4375,"y":151},{"x":594.4375,"y":115},{"x":636.4375,"y":114},{"x":628.4375,"y":61},{"x":583.4375,"y":66},{"x":517.4375,"y":65},{"x":521.4375,"y":466},{"x":582.4375,"y":465},{"x":583.4375,"y":512},{"x":718.4375,"y":513},{"x":722.4375,"y":583},{"x":525.4375,"y":582},{"x":450.4375,"y":584},{"x":448.4375,"y":466},{"x":498.4375,"y":469},{"x":491.4375,"y":64},{"x":434.4375,"y":64},{"x":448.4375,"y":427},{"x":473.4375,"y":428},{"x":463.4375,"y":90},{"x":480.4375,"y":88},{"x":484.4375,"y":446},{"x":412.4375,"y":447},{"x":366.4375,"y":580},{"x":202.4375,"y":582},{"x":48.4375,"y":582}],[{"x":141.4375,"y":68},{"x":144.4375,"y":189},{"x":47.4375,"y":186},{"x":44.4375,"y":69}],[{"x":368.4375,"y":65},{"x":370.4375,"y":190},{"x":265.4375,"y":184},{"x":262.4375,"y":68}],[{"x":374.4375,"y":268},{"x":378.4375,"y":374},{"x":274.4375,"y":374},{"x":268.4375,"y":279}],[{"x":153.4375,"y":281},{"x":159.4375,"y":380},{"x":55.4375,"y":377},{"x":51.4375,"y":280}],[{"x":241.4375,"y":207},{"x":242.4375,"y":259},{"x":191.4375,"y":261},{"x":188.4375,"y":209}],[{"x":216.4375,"y":109},{"x":219.4375,"y":135},{"x":191.4375,"y":136},{"x":191.4375,"y":113}],[{"x":232.4375,"y":327},{"x":233.4375,"y":358},{"x":198.4375,"y":359},{"x":196.4375,"y":331}],[{"x":51.4375,"y":432},{"x":372.4375,"y":429},{"x":314.4375,"y":543},{"x":25.4375,"y":546}],[{"x":387.4375,"y":470},{"x":353.4375,"y":551},{"x":335.4375,"y":541},{"x":378.4375,"y":460}],[{"x":376.4375,"y":392},{"x":370.4375,"y":418},{"x":278.4375,"y":417},{"x":281.4375,"y":401}]]}`
const floorPlan = `{"canvas":{"w":800,"h":600},"polygons":[[{"x":0,"y":0},{"x":120,"y":0},{"x":120,"y":340},{"x":180,"y":340},{"x":180,"y":-120},{"x":300,"y":-120},{"x":300,"y":340},{"x":360,"y":340},{"x":360,"y":0},{"x":480,"y":0},{"x":480,"y":420},{"x":300,"y":420},{"x":300,"y":540},{"x":340,"y":540},{"x":340,"y":720},{"x":140,"y":720},{"x":140,"y":540},{"x":180,"y":540},{"x":180,"y":420},{"x":0,"y":420}]]}`

func main() {
	polygon, holes, screen, err := polyjson.NewPolygonsFromJSON([]byte(floorPlan))
	if err != nil {
		panic(err)
	}

	var (
		ctx                = context.Background()
		start              = geom.Vector2{X: 30, Y: 30}
		dest               = geom.Vector2{X: 240, Y: 600}
		squareSize float32 = 16.55
		gridOffset         = geom.Vector2{X: 0.0, Y: 0.0}
		gridGraph          = grid.NewGrid(polygon, holes, squareSize, grid.WithOffset(gridOffset))
		pathfinder         = pathfind.NewPathfinder[geom.Vector2]([]graphs.NavGraph[geom.Vector2]{
			gridGraph,
		})
		dynamicObstacles = []obstacles.Obstacle{
			obstacles.GenerateCircle(geom.Vector2{X: 240, Y: 380}, 20, 15),
			obstacles.GenerateCircle(geom.Vector2{X: 148, Y: 382}, 18, 15),
			obstacles.GenerateCircle(geom.Vector2{X: 240, Y: 524}, 30, 30),
			obstacles.GenerateRectangle(geom.Vector2{X: 403, Y: 145}, 500, 20),
		}
		movingObstacles = []*MovingObstacle{
			newMovingObstacle(dynamicObstacles[0], moveHorizontal, 200, 1.5),
			newMovingObstacle(dynamicObstacles[1], moveVertical, 20, 0.1),
			newMovingObstacle(dynamicObstacles[2], moveDiagonal, 40, 0.2),
			newMovingObstacle(dynamicObstacles[3], moveHorizontal, 200, 0.2),
		}
		camera          = rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(-screen.X/2, -screen.Y/2), 0, 0.5)
		graphId         = 0
		isDrawGraph     = true
		isDrawSquares   = false
		lastSearchTime  = time.Now()
		edgesCount      = 0
		vertexCount     = 0
		squaresCount    = 0
		visSquaresCount = 0

		initTime string
		pathTime string
		path     []geom.Vector2 = nil
	)

	t := time.Now()
	if err := pathfinder.Initialize(ctx); err != nil {
		panic(err)
	}
	initTime = time.Since(t).String()

	visGraph := pathfinder.Graph(graphId)
	for _, edges := range visGraph {
		vertexCount++
		edgesCount += len(edges)
	}

	squaresCount = len(gridGraph.Squares())
	visSquaresCount = len(gridGraph.VisibleSquares())

	t2 := time.Now()
	path = pathfinder.Path(graphId, start, dest, pathfind.WithObstacles(dynamicObstacles))
	pathTime = time.Since(t2).String()

	rl.SetTraceLogLevel(rl.LogError)
	rl.SetTargetFPS(60)
	rl.InitWindow(int32(screen.X), int32(screen.Y), "Test *A")
	for {
		if rl.WindowShouldClose() {
			break
		}
		// regenerate graph and pathfinder based on sliders
		{
			gridGraph = grid.NewGrid(polygon, holes, float32(squareSize), grid.WithOffset(gridOffset))
			pathfinder = pathfind.NewPathfinder[geom.Vector2]([]graphs.NavGraph[geom.Vector2]{
				gridGraph,
			})

			pathfinder.Initialize(ctx)

			vertexCount = 0
			edgesCount = 0
			visGraph = pathfinder.Graph(graphId)
			for _, edges := range visGraph {
				vertexCount++
				edgesCount += len(edges)
			}

			squaresCount = len(gridGraph.Squares())
			visSquaresCount = len(gridGraph.VisibleSquares())
		}

		// moving dynamic obstacles
		for _, mo := range movingObstacles {
			mo.Move()
		}

		mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
		moveCamera(&camera, mouseWorldPos)

		rl.BeginDrawing()

		rl.BeginMode2D(camera)
		rl.ClearBackground(rl.White)

		// drawing map
		drawMap(polygon, holes, dynamicObstacles)
		if isDrawSquares {
			drawSquares(gridGraph.Squares())
		}
		if isDrawGraph {
			dGraph := pathfinder.GraphWithSearchPath(graphId, start, dest, pathfind.WithObstacles(dynamicObstacles))
			drawGraph(dGraph)
		}

		drawPath(start, dest, path, camera.Zoom, true)

		t := time.Now()
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			pp := geom.Vector2{X: mouseWorldPos.X, Y: mouseWorldPos.Y}
			if gridGraph.ContainsPoint(pp) {
				dest = pp
				t3 := time.Now()
				searchPath := pathfinder.Path(graphId, start, dest, pathfind.WithObstacles(dynamicObstacles))
				if len(searchPath) > 2 {
					path = searchPath
					pathTime = time.Since(t3).String()
				}
			}

			lastSearchTime = t
		}

		// recalc search path
		if time.Since(lastSearchTime) > searchPathTimeInterval {
			t3 := time.Now()
			searchPath := pathfinder.Path(graphId, start, dest, pathfind.WithObstacles(dynamicObstacles))
			if len(searchPath) > 2 {
				path = searchPath
				pathTime = time.Since(t3).String()
			} else {
				path = nil
			}

			lastSearchTime = t
		}

		rl.EndMode2D()

		drawTopPanel(int32(screen.X), rl.Vector2{}, &isDrawGraph, &isDrawSquares, initTime,
			pathTime, &squareSize, squaresCount, visSquaresCount, vertexCount, edgesCount, &gridOffset)
		rl.SetWindowTitle(fmt.Sprintf("Test A* (%v, %v)", int(mouseWorldPos.X), int(mouseWorldPos.Y)))

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func drawGraph(graph map[geom.Vector2][]geom.Vector2) {
	for p, elems := range graph {
		for _, elem := range elems {
			rl.DrawLine(int32(p.X), int32(p.Y), int32(elem.X), int32(elem.Y), rl.NewColor(230, 41, 55, 30))
		}
	}
}

func drawMap(polygon []geom.Vector2, holes [][]geom.Vector2, dynamicObstacles []obstacles.Obstacle) {
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

	for _, obstacle := range dynamicObstacles {
		polygon := obstacle.GetPolygon()

		n := len(polygon)
		if n < 2 {
			return
		}

		for i := range len(polygon) - 1 {
			p1, p2 := polygon[i], polygon[i+1]
			rl.DrawLine(int32(p1.X), int32(p1.Y), int32(p2.X), int32(p2.Y), rl.Magenta)
		}

		// DRAW END LINES
		startPos, endPos := polygon[n-1], polygon[0]
		rl.DrawLine(int32(startPos.X), int32(startPos.Y), int32(endPos.X), int32(endPos.Y), rl.Magenta)
	}
}

func drawSquares(squares []grid.Square) {
	for _, square := range squares {
		//rl.DrawCircle(int32(square.Center.X), int32(square.Center.Y), 1.0, rl.Blue)
		edges := square.Edges()
		for _, edge := range edges {
			rl.DrawLine(int32(edge.A.X), int32(edge.A.Y), int32(edge.B.X), int32(edge.B.Y), rl.NewColor(230, 41, 55, 20))
		}
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

func drawTopPanel(width int32, tPos rl.Vector2, isDrawGraph, isDrawSquares *bool, initTime,
	pathTime string, squareSize *float32, squaresCount, visSquaresCount, vertexCount, edgesCount int, offset *geom.Vector2) {

	*squareSize = rlgui.Slider(rl.Rectangle{600, 40, 120, 20}, "Squares",
		fmt.Sprintf("%.2f", *squareSize), *squareSize, 1, 50)

	(*offset).X = rlgui.Slider(rl.Rectangle{600, 70, 120, 20}, "OffsetX",
		fmt.Sprintf("%.2f", offset.X), offset.X, -30, 30)

	(*offset).Y = rlgui.Slider(rl.Rectangle{600, 100, 120, 20}, "OffsetY",
		fmt.Sprintf("%.2f", offset.Y), offset.Y, -30, 30)

	rl.DrawRectangle(int32(tPos.X), int32(tPos.Y), width, 30, rl.NewColor(127, 106, 79, 100))

	*isDrawGraph = rlgui.CheckBox(rl.NewRectangle(15, 10, 15, 15), "draw graph", *isDrawGraph)
	*isDrawSquares = rlgui.CheckBox(rl.NewRectangle(100, 10, 15, 15), "draw squares", *isDrawSquares)

	rl.DrawText(" | ", 200, 10, 15, rl.Gray)
	rlgui.Label(rl.NewRectangle(220, 10, 150, 15), "Init time: "+initTime)
	rl.DrawText(" | ", 320, 10, 15, rl.Gray)
	rlgui.Label(rl.NewRectangle(340, 10, 150, 15), "Path time: "+pathTime)
	rl.DrawText(" | ", 440, 10, 15, rl.Gray)
	rlgui.Label(rl.NewRectangle(460, 10, 150, 15), fmt.Sprintf("Square size: %v", *squareSize))
	rl.DrawText(" | ", 540, 10, 15, rl.Gray)
	rlgui.Label(rl.NewRectangle(560, 10, 200, 15), fmt.Sprintf("Squares %v/%v", visSquaresCount, squaresCount))
	rl.DrawText(" | ", 665, 10, 15, rl.Gray)
	rlgui.Label(rl.NewRectangle(680, 10, 180, 15), fmt.Sprintf("Graph (%vx%v)", vertexCount, edgesCount))
}

type MoveType int

const (
	moveHorizontal MoveType = iota + 1
	moveVertical
	moveDiagonal
)

type MovingObstacle struct {
	startPos      geom.Vector2
	obstacle      obstacles.Obstacle
	moveType      MoveType
	maxDistance   float32
	speed         float32
	moveDirection int
}

func newMovingObstacle(obstacle obstacles.Obstacle, moveType MoveType, maxDistance, speed float32) *MovingObstacle {
	return &MovingObstacle{
		startPos:      obstacle.GetCenter(),
		obstacle:      obstacle,
		moveType:      moveType,
		maxDistance:   maxDistance,
		speed:         speed,
		moveDirection: 1,
	}
}

func (m *MovingObstacle) Move() {
	movingObstacleCenter := m.obstacle.GetCenter()
	if geom.Distance(movingObstacleCenter, m.startPos) > m.maxDistance {
		m.moveDirection *= -1
	}

	switch m.moveType {
	case moveHorizontal:
		m.obstacle.Move(geom.Vector2{X: float32(m.moveDirection) * m.speed, Y: 0})
	case moveVertical:
		m.obstacle.Move(geom.Vector2{X: 0, Y: float32(m.moveDirection) * m.speed})
	case moveDiagonal:
		m.obstacle.Move(geom.Vector2{X: float32(m.moveDirection) * m.speed, Y: float32(m.moveDirection) * m.speed})
	}

}
