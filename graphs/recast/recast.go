package recast

import (
	"context"
	"math"

	"github.com/bolom009/astar"
	"github.com/bolom009/delaunay"
	"github.com/bolom009/geom"
	goclipper2 "github.com/bolom009/go-clipper2"
	"github.com/bolom009/pathfind/graphs"
	"github.com/bolom009/pathfind/mesh"
)

type edge struct {
	a geom.Vector2
	b geom.Vector2
}

type Recast struct {
	polygons        []*mesh.Polygon
	edges           []*edge
	clippedPolygons []*mesh.Polygon
	raycasts        []*Raycast
	triangles       []Triangle
	visibilityGraph graphs.Graph[geom.Vector2]
	vertices        []geom.Vector2
	costFunc        astar.CostFunc[geom.Vector2]
	kdTree          *KDTree
	searchOutOfArea bool

	// extra obstacles
	obstaclePool                *obstaclePool
	extraClippedPolygons        []*mesh.Polygon
	extraEdges                  []*edge
	extraObstacleId2ClippedPoly map[obstaclePolyPair]goclipper2.PathsD
}

func NewRecast(polygons []*mesh.Polygon, options ...option) *Recast {
	r := &Recast{
		polygons:                    polygons,
		triangles:                   make([]Triangle, 0),
		visibilityGraph:             make(graphs.Graph[geom.Vector2]),
		raycasts:                    make([]*Raycast, len(polygons)),
		extraClippedPolygons:        nil,
		extraObstacleId2ClippedPoly: make(map[obstaclePolyPair]goclipper2.PathsD),
		obstaclePool:                newObstaclePool(30),
		costFunc:                    heuristicEvaluation,
	}

	for _, option := range options {
		option(r)
	}

	return r
}

func (r *Recast) Generate(_ context.Context) error {
	r.prepareRaycasts()

	oPolygons := make([]*mesh.Polygon, 0, len(r.polygons))
	triLen := 0
	k := 0
	// offset each polygon and their innerHole + obstacles
	// then union results from offsets for to get new sub polygons
	for _, polygon := range r.polygons {
		polyOffsets := r.getPolyOffsetsWithUnion(polygon)

		rHoles := make([]*mesh.Hole, 0, len(polygon.Holes()))
		subPolygons := make([][]geom.Vector2, 0)
		for _, polyOffset := range polyOffsets {
			cPoly := toPoint2(polyOffset)
			if !goclipper2.IsPositiveD(polyOffset) {
				rHoles = append(rHoles, mesh.NewObstacle(cPoly, 0, false))
			} else {
				subPolygons = append(subPolygons, cPoly)
			}
		}

		for _, subPolygon := range subPolygons {
			triLen += len(subPolygon)
			subPolygonHoles := make([]*mesh.Hole, 0, len(rHoles))
			for _, hole := range rHoles {
				if isPolygonAFullyInsideB(hole.Points(), subPolygon, true) {
					subPolygonHoles = append(subPolygonHoles, hole)
					triLen += len(hole.Points())
				}
			}

			oPolygons = append(oPolygons, mesh.NewPolygon(subPolygon, subPolygonHoles, nil, 0))
			k++
		}
	}

	// process recast data based on all polygons what we got during clipper2 process
	triangles := make([]Triangle, 0, triLen)
	for _, polygon := range oPolygons {
		innerHoles := make([]*mesh.Hole, len(polygon.InnerHoles()))
		for i, innerHole := range polygon.InnerHoles() {
			innerHoles[i] = mesh.NewObstacle(innerHole.Points(), 0, innerHole.Viewable())
		}

		obstacles := make([]*mesh.Hole, len(polygon.Obstacles()))
		for i, obstacle := range polygon.Obstacles() {
			obstacles[i] = mesh.NewObstacle(obstacle.Points(), 0, obstacle.Viewable())
		}

		poly := mesh.NewPolygon(polygon.Points(), innerHoles, obstacles, 0)
		r.clippedPolygons = append(r.clippedPolygons, poly)

		pp := make([]*delaunay.Point, 0)
		for _, p := range poly.Points() {
			pp = append(pp, delaunay.NewPoint(p.X, p.Y))
		}

		hh := make([][]*delaunay.Point, 0)
		for _, hole := range poly.Holes() {
			newHole := make([]*delaunay.Point, len(hole.Points()))
			for i, p := range hole.Points() {
				newHole[i] = delaunay.NewPoint(p.X, p.Y)
			}

			hh = append(hh, newHole)
		}

		d := delaunay.NewSweepContext(pp, hh)

		cTriangles := convertTriangles(d.Triangulate())
		triangles = append(triangles, cTriangles...)
	}

	r.extraClippedPolygons = make([]*mesh.Polygon, len(r.clippedPolygons))
	copy(r.extraClippedPolygons, r.clippedPolygons)

	r.prepareEdges(triLen)

	r.triangles = triangles
	// generate graph based on triangles
	r.visibilityGraph = r.generateGraph()
	r.vertices = make([]geom.Vector2, len(r.visibilityGraph))

	i := 0
	for v := range r.visibilityGraph {
		r.vertices[i] = v
		i++
	}

	//r.kdTree = BuildKDTree(r.vertices, 0)
	return nil
}

func (r *Recast) AggregationGraph(start, dest geom.Vector2, _ *graphs.NavOpts) graphs.Graph[geom.Vector2] {
	for _, polygon := range r.extraClippedPolygons {
		polyPoints := polygon.Points()
		polyHoles := polygon.Holes()
		if !isInsidePolygonWithHoles(polyPoints, polyHoles, start) {
			continue
		}

		if !isInsidePolygonWithHoles(polyPoints, polyHoles, dest) {
			continue
		}

		// if start and dest in same visibility we can return start-dest graph
		if isLineSegmentInsidePolygonOrHoles(polyPoints, polyHoles, start, dest) {
			return graphs.Graph[geom.Vector2]{
				start: []geom.Vector2{dest},
				dest:  []geom.Vector2{start},
			}
		}
	}

	var (
		vis     = r.visibilityGraph.Copy()
		startOk = false
		destOk  = false
	)
	for _, triangle := range r.triangles {
		p0 := triangle[0]
		p1 := triangle[1]
		p2 := triangle[2]

		if pointInsideTriangle(p0, p1, p2, start) {
			vis.LinkBoth(start, p0)
			vis.LinkBoth(start, p1)
			vis.LinkBoth(start, p2)
			startOk = true
			continue
		}

		if pointInsideTriangle(p0, p1, p2, dest) {
			vis.LinkBoth(p0, dest)
			vis.LinkBoth(p1, dest)
			vis.LinkBoth(p2, dest)
			destOk = true
		}
	}

	if r.searchOutOfArea {
		if !startOk {
			closestPoint, ok := r.closestPointOnPolygon(start)
			if ok {
				vis.LinkBoth(start, closestPoint)
				visiblePoints := r.getVisiblePoints(closestPoint)
				for _, visiblePoint := range visiblePoints {
					vis.LinkBoth(visiblePoint, closestPoint)
				}
			}
		}

		if !destOk {
			closestPoint, ok := r.closestPointOnPolygon(dest)
			if ok {
				vis.LinkBoth(dest, closestPoint)
				visiblePoints := r.getVisiblePoints(closestPoint)
				for _, visiblePoint := range visiblePoints {
					vis.LinkBoth(visiblePoint, closestPoint)
				}
			}
		}
	}

	return vis
}

func (r *Recast) AddObstacles(obstacles ...*mesh.Hole) []uint32 {
	if len(obstacles) == 0 {
		return nil
	}

	ids := make([]uint32, len(obstacles))
	for i, obstacle := range obstacles {
		rPoints := goclipper2.ReversePath(obstacle.Points())
		newPaths := inflatePathsD(goclipper2.PathsD{toPathD(rPoints)}, float64(obstacle.Offset()), goclipper2.Miter, goclipper2.Polygon)

		ids[i] = r.obstaclePool.New(obstacle)

	loop:
		for j, cPolygon := range r.clippedPolygons {
			polyPoints := cPolygon.Points()
			for _, p := range obstacle.Points() {
				if pointInPolygon(p, polyPoints) {
					key := obstaclePolyPair{oId: ids[i], pId: j}
					r.extraObstacleId2ClippedPoly[key] = newPaths

					break loop
				}
			}
		}
	}

	r.rebuild()
	return ids
}

func (r *Recast) RemoveObstacles(ids ...uint32) {
	for _, id := range ids {
		for opKey := range r.extraObstacleId2ClippedPoly {
			if opKey.oId == id {
				delete(r.extraObstacleId2ClippedPoly, opKey)
			}
		}

		r.obstaclePool.Delete(id)
	}

	r.rebuild()
}

func (r *Recast) rebuild() {
	oPolygons := make([]*mesh.Polygon, 0, len(r.clippedPolygons))
	for i, cPolygon := range r.clippedPolygons {
		extraClippedObstacles := make([]goclipper2.PathsD, 0)
		for opKey, pathsD := range r.extraObstacleId2ClippedPoly {
			if opKey.pId == i {
				extraClippedObstacles = append(extraClippedObstacles, pathsD)
			}
		}

		if len(extraClippedObstacles) == 0 {
			oPolygons = append(oPolygons, cPolygon)
			continue
		}

		subject := goclipper2.PathsD{toPathD(cPolygon.Points())}
		for _, hole := range cPolygon.Holes() {
			subject = append(subject, toPathD(hole.Points()))
		}

		for _, extraObstacle := range extraClippedObstacles {
			subject = append(subject, extraObstacle...)
		}

		// union clipped polygons with external subject
		polyOffsets := unionPathsD(subject, goclipper2.Positive)

		rHoles := make([]*mesh.Hole, 0, len(cPolygon.Holes()))
		subPolygons := make([][]geom.Vector2, 0)
		for _, polyOffset := range polyOffsets {
			cPoly := toPoint2(polyOffset)
			if !goclipper2.IsPositiveD(polyOffset) {
				rHoles = append(rHoles, mesh.NewObstacle(cPoly, 0, false))
			} else {
				subPolygons = append(subPolygons, cPoly)
			}
		}

		for _, subPolygon := range subPolygons {
			subPolygonHoles := make([]*mesh.Hole, 0, len(rHoles))
			for _, hole := range rHoles {
				if isPolygonAFullyInsideB(hole.Points(), subPolygon, true) {
					subPolygonHoles = append(subPolygonHoles, hole)
				}
			}

			oPolygons = append(oPolygons, mesh.NewPolygon(subPolygon, subPolygonHoles, nil, 0))
		}
	}

	// TODO do we need to recalc raycast? raycast - check visibility to enemy through obstacle
	r.extraClippedPolygons = r.extraClippedPolygons[:0]

	triangles := make([]Triangle, 0)
	for _, polygon := range oPolygons {
		innerHoles := make([]*mesh.Hole, len(polygon.InnerHoles()))
		for i, innerHole := range polygon.InnerHoles() {
			innerHoles[i] = mesh.NewObstacle(innerHole.Points(), 0, innerHole.Viewable())
		}

		oObstacles := make([]*mesh.Hole, len(polygon.Obstacles()))
		for i, obstacle := range polygon.Obstacles() {
			oObstacles[i] = mesh.NewObstacle(obstacle.Points(), 0, obstacle.Viewable())
		}

		poly := mesh.NewPolygon(polygon.Points(), innerHoles, oObstacles, 0)
		r.extraClippedPolygons = append(r.extraClippedPolygons, poly)

		pp := make([]*delaunay.Point, 0)
		for _, p := range poly.Points() {
			pp = append(pp, delaunay.NewPoint(p.X, p.Y))
		}

		hh := make([][]*delaunay.Point, 0)
		for _, hole := range poly.Holes() {
			newHole := make([]*delaunay.Point, len(hole.Points()))
			for i, p := range hole.Points() {
				newHole[i] = delaunay.NewPoint(p.X, p.Y)
			}

			hh = append(hh, newHole)
		}

		d := delaunay.NewSweepContext(pp, hh)

		cTriangles := convertTriangles(d.Triangulate())
		triangles = append(triangles, cTriangles...)
	}

	r.prepareEdges(len(r.edges))

	r.triangles = triangles
	// generate graph based on triangles
	r.visibilityGraph = r.generateGraph()
	r.vertices = make([]geom.Vector2, len(r.visibilityGraph))

	i := 0
	for v := range r.visibilityGraph {
		r.vertices[i] = v
		i++
	}
}

func (r *Recast) GetClosestPoint(point geom.Vector2) (geom.Vector2, bool) {
	// fastest way (~50ns) but return only vertex of polygon not closest point to polygon
	//closestPoint, _ := r.kdTree.search(point)
	//return closestPoint

	// ~1000ns
	closestPoint, ok := r.closestPointOnPolygon(point)
	return closestPoint, ok
}

func (r *Recast) IsRaycastHit(start, end geom.Vector2) bool {
	for _, raycast := range r.raycasts {
		if raycast.CheckLineIntersectsPolygon(start, end) {
			return true
		}
	}

	return false
}

func (r *Recast) GetVisibility(_ *graphs.NavOpts) graphs.Graph[geom.Vector2] {
	return r.visibilityGraph.Copy()
}

func (r *Recast) ContainsPoint(point geom.Vector2) bool {
	//for _, triangle := range r.triangles {
	//	p0 := triangle[0]
	//	p1 := triangle[1]
	//	p2 := triangle[2]
	//
	//	if pointInsideTriangle(p0, p1, p2, point) {
	//		return true
	//	}
	//}
	//
	//return false

	for _, polygon := range r.polygons {
		if isInsidePolygonWithHoles(polygon.Points(), polygon.Holes(), point) {
			return true
		}
	}

	return false
}

func (r *Recast) Cost(a geom.Vector2, b geom.Vector2) float32 {
	return r.costFunc(a, b)
}

const (
	scaleFactor        = 1e6
	hashRnd            = 1099511628211
	hashDefault uint64 = 14695981039346656037
)

func (r *Recast) HashIndex(v geom.Vector2) int64 {
	qx := quantizeFloat(v.X)
	qy := quantizeFloat(v.Y)

	var hash = hashDefault
	hash = (hash * hashRnd) ^ uint64(qx)
	hash = (hash * hashRnd) ^ uint64(qy)

	return int64(hash)
}

func quantizeFloat(f float32) int64 {
	return int64(f * scaleFactor)
}

func (r *Recast) Triangles() []Triangle {
	return r.triangles
}

func (r *Recast) generateGraph() graphs.Graph[geom.Vector2] {
	vis := make(graphs.Graph[geom.Vector2], len(r.triangles))
	for _, triangle := range r.triangles {
		var (
			p0 = triangle[0]
			p1 = triangle[1]
			p2 = triangle[2]
		)

		// link top vertices
		vis.LinkBoth(p0, p1, 2)
		vis.LinkBoth(p0, p2, 2)
		vis.LinkBoth(p1, p2, 2)
	}

	return vis
}

func (r *Recast) prepareEdges(edgeLen int) {
	r.edges = make([]*edge, 0, edgeLen)
	for _, polygon := range r.polygons {
		r.edges = append(r.edges, polyToEdges(polygon.Points())...)
		for _, hole := range polygon.Holes() {
			r.edges = append(r.edges, polyToEdges(hole.Points())...)
		}
	}

	extraObstacles := r.obstaclePool.GetList()
	for _, extraObstacle := range extraObstacles {
		r.edges = append(r.edges, polyToEdges(extraObstacle.Points())...)
	}

	r.extraEdges = make([]*edge, 0, edgeLen)
	for _, polygon := range r.extraClippedPolygons {
		r.extraEdges = append(r.extraEdges, polyToEdges(polygon.Points())...)
		for _, hole := range polygon.Holes() {
			r.extraEdges = append(r.extraEdges, polyToEdges(hole.Points())...)
		}
	}
}

func polyToEdges(points []geom.Vector2) []*edge {
	pLen := len(points)
	edges := make([]*edge, pLen)
	for i := range pLen - 1 {
		p1, p2 := points[i], points[i+1]
		edges[i] = &edge{p1, p2}
	}

	edges[pLen-1] = &edge{points[0], points[pLen-1]}

	return edges
}

func (r *Recast) prepareRaycasts() {
	for i, polygon := range r.polygons {
		r.raycasts[i] = NewRaycast(polygon.Points(), polygon.Obstacles())
	}
}

func isInsidePolygonWithHoles(polygon []geom.Vector2, holes []*mesh.Hole, point geom.Vector2) bool {
	if !pointInPolygon(point, polygon) {
		return false
	}
	for _, hole := range holes {
		if pointInPolygon(point, hole.Points()) {
			return false
		}
	}
	return true
}

func isLineSegmentInsidePolygonOrHoles(polygon []geom.Vector2, holes []*mesh.Hole, lineStart, lineEnd geom.Vector2) bool {
	// Check if the line segment intersects any of the polygon edges
	if isLineSegmentInsidePolygon(polygon, lineStart, lineEnd) {
		return false
	}

	for _, hole := range holes {
		points := hole.Points()
		for i := 0; i < len(points)-1; i++ {
			if doLinesIntersect(lineStart, lineEnd, points[i], points[i+1]) {
				return false // Intersects a hole
			}
		}
	}

	return true // Line segment is completely inside the polygon
}

// pointInPolygon checks if a point p is inside a polygon using the ray casting method.
func pointInPolygon(p geom.Vector2, poly []geom.Vector2) bool {
	inside := false
	n := len(poly)
	for i := 0; i < n; i++ {
		p1 := poly[i]
		p2 := poly[(i+1)%n]

		// Check if the edge (p1->p2) straddles the horizontal line at p.Y
		condY := (p1.Y <= p.Y && p2.Y > p.Y) || (p2.Y <= p.Y && p1.Y > p.Y)
		if condY {
			// Compute the x-coordinate of intersection of the polygon edge with the line y = p.Y
			xIntersect := p1.X + (p.Y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y)
			if xIntersect > p.X {
				inside = !inside
			}
		}
	}
	return inside
}

// isLineSegmentInsidePolygon function to check if a line segment is inside a polygon
func isLineSegmentInsidePolygon(polygon []geom.Vector2, lineStart, lineEnd geom.Vector2) bool {
	n := len(polygon)
	for i := 0; i < n; i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%n]
		if doLinesIntersect(lineStart, lineEnd, p1, p2) {
			return true // The line segment intersects the polygon edge
		}
	}

	return false
}

func doLinesIntersect(p1, p2, q1, q2 geom.Vector2) bool {
	// Convert points to segments
	// Calculate the orientation
	o1 := orientation(p1, p2, q1)
	o2 := orientation(p1, p2, q2)
	o3 := orientation(q1, q2, p1)
	o4 := orientation(q1, q2, p2)

	// General case
	if o1 != o2 && o3 != o4 {
		return true
	}

	// Special cases (collinear points)
	if o1 == 0 && onSegment(p1, q1, p2) {
		return true
	}
	if o2 == 0 && onSegment(p1, q2, p2) {
		return true
	}
	if o3 == 0 && onSegment(q1, p1, q2) {
		return true
	}
	if o4 == 0 && onSegment(q1, p2, q2) {
		return true
	}
	return false
}

// orientation helper function for intersection calculations
func orientation(p, q, r geom.Vector2) int {
	val := (q.Y-r.Y)*(p.X-q.X) - (q.X-r.X)*(p.Y-q.Y)
	if val == 0 {
		return 0 // Collinear
	}
	if val > 0 {
		return 1 // Clockwise
	}

	return 2 // Counterclockwise
}

// onSegment helper function for segment calculations
func onSegment(p, q, r geom.Vector2) bool {
	var (
		pX, pY = float64(p.X), float64(p.Y)
		qX, qY = float64(q.X), float64(q.Y)
		rX, rY = float64(r.X), float64(r.Y)
	)

	return qX <= math.Max(pX, rX) && qX >= math.Min(pX, rX) && qY <= math.Max(pY, rY) && qY >= math.Min(pY, rY)
}

// pointOnSegment helper function for check point on segment with epsilon
func pointOnSegment(p, a, b geom.Vector2, eps float64) bool {
	var (
		pX, pY = float64(p.X), float64(p.Y)
		aX, aY = float64(a.X), float64(a.Y)
		bX, bY = float64(b.X), float64(b.Y)
	)

	if pX < math.Min(aX, bX)-eps || pX > math.Max(aX, bX)+eps ||
		pY < math.Min(aY, bY)-eps || pY > math.Max(aY, bY)+eps {
		return false
	}

	cross := (b.X-a.X)*(p.Y-a.Y) - (b.Y-a.Y)*(p.X-a.X)
	return math.Abs(float64(cross)) <= eps
}

func (r *Recast) getVisiblePoints(point geom.Vector2) []geom.Vector2 {
	visiblePoints := make([]geom.Vector2, len(r.extraEdges))
	count := 0

	// TODO very fast check ~11146ns
	for _, edge := range r.extraEdges {
		if pointOnSegment(point, edge.a, edge.b, 1e-1) {
			visiblePoints[count] = edge.a
			count++
			visiblePoints[count] = edge.b
			count++
		}
	}

	// TODO better but not so accurate :( ~99853ns
	//for _, v := range r.vertices {
	//	hasIntersection := false
	//	for _, edge := range r.edges {
	//		if lineSegmentIntersection(edge.a, edge.b, v, point) {
	//			hasIntersection = true
	//			break
	//		}
	//	}
	//
	//	if !hasIntersection {
	//		visiblePoints[count] = v
	//		count++
	//	}
	//}

	// TODO more accurate but expensive :( ~190585ns
	//for _, polygon := range r.polygons {
	//	points := polygon.Points()
	//	if !pointInPolygon(point, points) {
	//		continue
	//	}
	//
	//	holes := polygon.Holes()
	//	for _, v := range r.vertices {
	//		if isLineSegmentInsidePolygonOrHoles(points, holes, point, v) {
	//			visiblePoints[count] = v
	//			count++
	//		}
	//	}
	//}

	return visiblePoints[:count]
}

// closestPointOnPolygon finds the closest point on the polygon's boundary to the given point
func (r *Recast) closestPointOnPolygon(startPoint geom.Vector2) (geom.Vector2, bool) {
	minDist := float32(math.MaxFloat32)
	closestFoundPoint := geom.Vector2{}
	exist := false

	// Check exterior ring
	for _, polygon := range r.extraClippedPolygons {
		clippedPolygon := polygon.Points()
		for i := 0; i < len(clippedPolygon); i++ {
			p1 := clippedPolygon[i]
			p2 := clippedPolygon[(i+1)%len(clippedPolygon)] // Wrap around for last segment

			closestOnSeg := closestPointOnSegment(startPoint, p1, p2)
			dist := geom.Distance(startPoint, closestOnSeg)

			if dist < minDist {
				minDist = dist
				closestFoundPoint = closestOnSeg
				exist = true
			}
		}

		// Check interior rings (holes)
		for _, hole := range polygon.Holes() {
			holePoints := hole.Points()
			for i := 0; i < len(holePoints); i++ {
				p1 := holePoints[i]
				p2 := holePoints[(i+1)%len(holePoints)]

				closestOnSeg := closestPointOnSegment(startPoint, p1, p2)
				dist := geom.Distance(startPoint, closestOnSeg)

				if dist < minDist {
					minDist = dist
					closestFoundPoint = closestOnSeg
					exist = true
				}
			}
		}
	}

	return closestFoundPoint, exist
}

func closestPointOnSegment(p, a, b geom.Vector2) geom.Vector2 {
	ap := geom.Vector2{X: p.X - a.X, Y: p.Y - a.Y}
	ab := geom.Vector2{X: b.X - a.X, Y: b.Y - a.Y}

	dotProd := ap.X*ab.X + ap.Y*ab.Y
	lenSq := ab.X*ab.X + ab.Y*ab.Y

	if lenSq == 0 { // a and b are the same point
		return a
	}

	t := dotProd / lenSq
	if t < 0 {
		return a
	} else if t > 1 {
		return b
	}

	return geom.Vector2{X: a.X + t*ab.X, Y: a.Y + t*ab.Y}
}

func heuristicEvaluation(a, b geom.Vector2) float32 {
	x := a.X - b.X
	y := a.Y - b.Y

	return float32(math.Sqrt(float64(x*x + y*y)))
}
