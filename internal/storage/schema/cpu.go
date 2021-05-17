package schema

import (
	"sync"
	"time"
)

type CPULoadEntity struct {
	timestamp  int64
	userMode   float32
	systemMode float32
	idle       float32
}

type CPULoadDto struct {
	UserMode   float32
	SystemMode float32
	Idle       float32
}

type CPULoadTable struct {
	entities []*CPULoadEntity
	mtx      *sync.RWMutex
}

func (t *CPULoadTable) Init() *CPULoadTable {
	t.entities = []*CPULoadEntity{}
	t.mtx = &sync.RWMutex{}
	return t
}

func (t *CPULoadTable) AddEntity(d CPULoadDto) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	e := &CPULoadEntity{
		timestamp:  time.Now().Unix(),
		userMode:   d.UserMode,
		systemMode: d.SystemMode,
		idle:       d.Idle,
	}
	t.entities = append(t.entities, e)
}

func (t *CPULoadTable) GetAverage(period int32) CPULoadDto {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	currentTime := time.Now().Unix()
	sum := &CPULoadDto{}
	num := 0
	for i := len(t.entities) - 1; i >= 0; i-- {
		if t.entities[i].timestamp < currentTime-int64(period) {
			break
		}
		num++
		sum.Idle += t.entities[i].idle
		sum.SystemMode += t.entities[i].systemMode
		sum.UserMode += t.entities[i].userMode
	}

	return CPULoadDto{
		UserMode:   sum.UserMode / float32(num),
		SystemMode: sum.SystemMode / float32(num),
		Idle:       sum.Idle / float32(num),
	}
}
