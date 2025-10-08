package recast

import (
	"testing"

	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind/mesh"
	"github.com/stretchr/testify/assert"
)

func TestObstaclePool(t *testing.T) {
	obstacle := []geom.Vector2{{1, 1}, {1, 3}, {3, 1}}
	obstacle2 := []geom.Vector2{{10, 10}, {10, 30}, {30, 10}}
	deleteList := make([]uint32, 0)

	pool := NewObstaclePool(30)
	idx := pool.New(mesh.NewObstacle(obstacle, 3, true))
	deleteList = append(deleteList, idx)
	assert.Equal(t, uint32(1), idx)
	assert.Equal(t, 1, len(pool.GetList()))

	idx = pool.New(mesh.NewObstacle(obstacle, 3, true))
	deleteList = append(deleteList, idx)
	assert.Equal(t, uint32(2), idx)
	assert.Equal(t, 2, len(pool.GetList()))

	pool.Delete(deleteList[1])

	idx = pool.New(mesh.NewObstacle(obstacle2, 3, true))
	deleteList = append(deleteList, idx)
	assert.Equal(t, uint32(2), idx)
	assert.Equal(t, 2, len(pool.GetList()))

	pool.Delete(deleteList[2])
	assert.Equal(t, 1, len(pool.GetList()))

	pool.Delete(deleteList[0])
	assert.Equal(t, 0, len(pool.GetList()))

	idx = pool.New(mesh.NewObstacle(obstacle2, 3, true))
	assert.Equal(t, uint32(1), idx)
}
