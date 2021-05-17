package schema

import (
	"sync"
	"time"
)

type Protocol string

const (
	TCP Protocol = "TCP"
	UDP          = "UDP"
)

type SocketLoadEntity struct {
	timestamp int64
	info      SocketLoadInfo
}

type SocketLoadDto struct {
	command  string
	pid      int32
	user     string
	protocol Protocol
	port     int32
}

type SocketLoadInfo struct {
	items []SocketLoadDto
}

type SocketLoadTable struct {
	entities []*SocketLoadEntity
	mtx      *sync.RWMutex
}

func (t *SocketLoadTable) Init() *SocketLoadTable {
	t.entities = []*SocketLoadEntity{}
	t.mtx = &sync.RWMutex{}
	return t
}

func (t *SocketLoadTable) AddEntity(d SocketLoadInfo) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	e := &SocketLoadEntity{
		timestamp: time.Now().Unix(),
		info:      d,
	}
	t.entities = append(t.entities, e)
}

func (t *SocketLoadTable) GetAverage(period int32) SocketLoadInfo {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	currentTime := time.Now().Unix()
	res := map[SocketLoadDto]bool{}
	for i := len(t.entities) - 1; i >= 0; i-- {
		if t.entities[i].timestamp < currentTime-int64(period) {
			break
		}
		for _, v := range t.entities[i].info.items {
			if _, ok := res[v]; !ok {
				res[v] = true
			}
		}
	}

	arr := make([]SocketLoadDto, 0, len(res))
	for k := range res {
		arr = append(arr, k)
	}

	return SocketLoadInfo{
		items: arr,
	}
}
