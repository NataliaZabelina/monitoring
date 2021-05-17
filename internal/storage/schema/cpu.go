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
	User_mode   float32
	System_mode float32
	Idle        float32
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
		user_mode:   d.User_mode,
		system_mode: d.System_mode,
		idle:        d.Idle,
	}
	t.entities = append(t.entities, e)
}

func (t *CpuLoadTable) GetAverage(period int32) CpuLoadDto {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	currentTime := time.Now().Unix()
	sum := &CpuLoadDto{}
	num := 0
	for i := len(t.entities) - 1; i >= 0; i-- {
		if t.entities[i].timestamp < currentTime-int64(period) {
			break
		}
		num++
		sum.Idle += t.entities[i].idle
		sum.System_mode += t.entities[i].system_mode
		sum.User_mode += t.entities[i].user_mode
	}

	return CpuLoadDto{
		User_mode:   sum.User_mode / float32(num),
		System_mode: sum.System_mode / float32(num),
		Idle:        sum.Idle / float32(num),
	}
}
