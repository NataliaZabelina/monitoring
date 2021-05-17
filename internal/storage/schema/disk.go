package schema

import (
	"sync"
	"time"
)

type DiskIO struct {
	tps float32
	kbs float32
	cpu float32
}

type DiskFS struct {
	mused float32
	iused float32
}

type DiskLoadEntity struct {
	timestamp int64
	disk_io   map[string]DiskIO
	disk_fs   map[string]DiskFS
}

type DiskLoadDto struct {
	disk_io map[string]DiskIO
	disk_fs map[string]DiskFS
}

type DiskLoadTable struct {
	entities []*DiskLoadEntity
	mtx      *sync.RWMutex
}

func (t *DiskLoadTable) Init() *DiskLoadTable {
	t.entities = []*DiskLoadEntity{}
	t.mtx = &sync.RWMutex{}
	return t
}

func (t *DiskLoadTable) AddEntity(d DiskLoadDto) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	e := &DiskLoadEntity{
		timestamp: time.Now().Unix(),
		disk_io:   d.disk_io,
		disk_fs:   d.disk_fs,
	}
	t.entities = append(t.entities, e)
}

func (t *DiskLoadTable) GetAverage(period int32) DiskLoadDto {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	currentTime := time.Now().Unix()
	sum := &DiskLoadDto{
		disk_io: map[string]DiskIO{},
		disk_fs: map[string]DiskFS{},
	}
	num := 0
	for i := len(t.entities) - 1; i >= 0; i-- {
		if t.entities[i].timestamp < currentTime-int64(period) {
			break
		}
		num++
		for k, v := range t.entities[i].disk_fs {
			if _, ok := sum.disk_fs[k]; !ok {
				sum.disk_fs[k] = v
				continue
			}
			sum.disk_fs[k] = DiskFS{
				mused: sum.disk_fs[k].mused + v.mused,
				iused: sum.disk_fs[k].iused + v.iused,
			}
		}

		for k, v := range t.entities[i].disk_io {
			if _, ok := sum.disk_io[k]; !ok {
				sum.disk_io[k] = v
				continue
			}
			sum.disk_io[k] = DiskIO{
				tps: sum.disk_io[k].tps + v.tps,
				kbs: sum.disk_io[k].kbs + v.kbs,
				cpu: sum.disk_io[k].cpu + v.cpu,
			}
		}
	}

	result := DiskLoadDto{
		disk_io: map[string]DiskIO{},
		disk_fs: map[string]DiskFS{},
	}

	for k, v := range sum.disk_fs {
		result.disk_fs[k] = DiskFS{
			mused: v.mused/float32(num),
			iused: v.iused/float32(num),
		}
	}

	for k, v := range sum.disk_io {
		result.disk_io[k] = DiskIO{
			tps: v.tps/float32(num),
			kbs: v.kbs/float32(num),
			cpu: v.cpu/float32(num),
		}
	}

	return result
}
