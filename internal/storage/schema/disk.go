package schema

import (
	"sync"
	"time"
)

type Disk struct {
	Param1 float32
	Param2 float32
}

type DiskEntity struct {
	timestamp int64
	Disk      map[string]Disk
}

type DiskDto struct {
	Disk map[string]Disk
}

type DiskTable struct {
	entities []*DiskEntity
	mtx      *sync.RWMutex
}

func (t *DiskTable) Init() *DiskTable {
	t.entities = []*DiskEntity{}
	t.mtx = &sync.RWMutex{}
	return t
}

func (t *DiskTable) AddEntity(d DiskDto) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	e := &DiskEntity{
		timestamp: time.Now().Unix(),
		Disk:      d.Disk,
	}
	t.entities = append(t.entities, e)
}

func (t *DiskTable) GetAverage(period int32) DiskDto {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	currentTime := time.Now().Unix()
	sum := &DiskDto{
		Disk: map[string]Disk{},
	}
	num := 0
	for i := len(t.entities) - 1; i >= 0; i-- {
		if t.entities[i].timestamp < currentTime-int64(period) {
			break
		}
		num++
		for k, v := range t.entities[i].Disk {
			if _, ok := sum.Disk[k]; !ok {
				sum.Disk[k] = v
				continue
			}
			sum.Disk[k] = Disk{
				Param2: sum.Disk[k].Param2 + v.Param2,
				Param1: sum.Disk[k].Param1 + v.Param1,
			}
		}
	}

	result := DiskDto{
		Disk: map[string]Disk{},
	}

	for k, v := range sum.Disk {
		result.Disk[k] = Disk{
			Param2: v.Param2 / float32(num),
			Param1: v.Param1 / float32(num),
		}
	}

	return result
}
