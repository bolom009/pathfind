package recast

import (
	"github.com/bolom009/pathfind/mesh"
)

type obstaclePolyPair struct {
	oId uint32
	pId int
}

type obstaclePool struct {
	byID   map[uint32]uint32
	items  []*mesh.Hole
	ids    []uint32
	free   []uint32
	nextID uint32
}

func newObstaclePool(cap int) *obstaclePool {
	return &obstaclePool{
		byID:   make(map[uint32]uint32, cap),
		items:  make([]*mesh.Hole, 0, cap),
		ids:    make([]uint32, 0, cap),
		free:   make([]uint32, 0, cap),
		nextID: 1,
	}
}

func (p *obstaclePool) New(h *mesh.Hole) uint32 {
	var id uint32
	if n := len(p.free); n > 0 {
		id = p.free[n-1]
		p.free = p.free[:n-1]
	} else {
		id = p.nextID
		p.nextID++
		if p.nextID == 0 {
			p.nextID = 1
		}
	}
	idx := uint32(len(p.items))
	p.items = append(p.items, h)
	p.ids = append(p.ids, id)
	p.byID[id] = idx
	return id
}

func (p *obstaclePool) Delete(id uint32) {
	idx, ok := p.byID[id]
	if !ok {
		return
	}

	last := uint32(len(p.items) - 1)
	if idx != last {
		p.items[idx] = p.items[last]
		movedID := p.ids[last]
		p.ids[idx] = movedID
		p.byID[movedID] = idx
	}

	p.items = p.items[:last]
	p.ids = p.ids[:last]
	delete(p.byID, id)
	p.free = append(p.free, id)
}

func (p *obstaclePool) GetList() []*mesh.Hole {
	return p.items
}
