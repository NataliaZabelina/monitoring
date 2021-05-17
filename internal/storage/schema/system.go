package schema

import (
	"sync"
	"time"
)

type SystemLoadEntity struct {
	timestamp int64
	loadAvg   float32
}

type SystemLoadDto struct {
	LoadAvg float32
}

type SystemLoadTable struct {
	entities []*SystemLoadEntity
	mtx      *sync.RWMutex
}

func (t *SystemLoadTable) Init() *SystemLoadTable {
	t.entities = []*SystemLoadEntity{}
	t.mtx = &sync.RWMutex{}
	return t
}

func (t *SystemLoadTable) AddEntity(d SystemLoadDto) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	e := &SystemLoadEntity{
		timestamp: time.Now().Unix(),
		loadAvg:   d.LoadAvg,
	}
	t.entities = append(t.entities, e)
}

func (t *SystemLoadTable) GetAverage(period int32) SystemLoadDto {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	currentTime := time.Now().Unix()
	sum := &SystemLoadDto{}
	num := 0
	for i := len(t.entities) - 1; i >= 0; i-- {
		if t.entities[i].timestamp < currentTime-int64(period) {
			break
		}
		num++
		sum.LoadAvg += t.entities[i].loadAvg
	}

	return SystemLoadDto{
		LoadAvg: sum.LoadAvg / float32(num),
	}
}
