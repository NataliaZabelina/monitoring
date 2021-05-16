package schema

import (
	"sync"
	"time"
)

type CpuLoadEntity struct {
	timestamp   int64
	user_mode   float32
	system_mode float32
	idle        float32
}

type CpuLoadDto struct {
	user_mode   float32
	system_mode float32
	idle        float32
}

type CpuLoadTable struct {
	entities []*CpuLoadEntity
	mtx      *sync.RWMutex
}

func (t *CpuLoadTable) Init() *CpuLoadTable {
	t.entities = []*CpuLoadEntity{}
	t.mtx = &sync.RWMutex{}
	return t
}

func (t *CpuLoadTable) AddEntity(d CpuLoadDto) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	e := &CpuLoadEntity{
		timestamp:   time.Now().Unix(),
		user_mode:   d.user_mode,
		system_mode: d.system_mode,
		idle:        d.idle,
	}
	t.entities = append(t.entities, e)
}

func (t *CpuLoadTable) GetAverage(period int32) CpuLoadDto {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	currentTime := time.Now().Unix()
	sum := &CpuLoadDto{}
	num := 0
	for i := len(t.entities) - 1; i < 0; i-- {
		if t.entities[i].timestamp < currentTime-int64(period) {
			break
		}
		num++
		sum.idle += t.entities[i].idle
		sum.system_mode += t.entities[i].system_mode
		sum.user_mode += t.entities[i].user_mode
	}

	return CpuLoadDto{
		user_mode:   sum.user_mode / float32(num),
		system_mode: sum.system_mode / float32(num),
		idle:        sum.idle / float32(num),
	}
}
