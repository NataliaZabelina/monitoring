package schema

import (
	"sync"
	"time"
)

type TCPCountEntity struct {
	timestamp int64
	tcp_data  map[string]int32
}

type TCPCountDto struct {
	tcp_data map[string]int32
}

type TCPCountTable struct {
	entities []*TCPCountEntity
	mtx      *sync.RWMutex
}

func (t *TCPCountTable) Init() *TCPCountTable {
	t.entities = []*TCPCountEntity{}
	t.mtx = &sync.RWMutex{}
	return t
}

func (t *TCPCountTable) AddEntity(d TCPCountDto) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	e := &TCPCountEntity{
		timestamp: time.Now().Unix(),
		tcp_data:  d.tcp_data,
	}
	t.entities = append(t.entities, e)
}

func (t *TCPCountTable) GetAverage(period int32) TCPCountDto {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	currentTime := time.Now().Unix()
	sum := &TCPCountDto{
		tcp_data: map[string]int32{},
	}
	num := int32(0)
	for i := len(t.entities) - 1; i >= 0; i-- {
		if t.entities[i].timestamp < currentTime-int64(period) {
			break
		}
		num++
		for k, v := range t.entities[i].tcp_data {
			if _, ok := sum.tcp_data[k]; !ok {
				sum.tcp_data[k] = v
				continue
			}
			sum.tcp_data[k] += v
		}
	}

	result := TCPCountDto{
		tcp_data: map[string]int32{},
	}

	for k, v := range sum.tcp_data {
		result.tcp_data[k] = v / num
	}

	return result
}
