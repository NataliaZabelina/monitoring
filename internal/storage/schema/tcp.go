package schema

import (
	"sync"
	"time"
)

type TCPCountEntity struct {
	timestamp int64
	tcpData   map[string]int32
}

type TCPCountDto struct {
	tcpData map[string]int32
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
		tcpData:   d.tcpData,
	}
	t.entities = append(t.entities, e)
}

func (t *TCPCountTable) GetAverage(period int32) TCPCountDto {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	currentTime := time.Now().Unix()
	sum := &TCPCountDto{
		tcpData: map[string]int32{},
	}
	num := int32(0)
	for i := len(t.entities) - 1; i >= 0; i-- {
		if t.entities[i].timestamp < currentTime-int64(period) {
			break
		}
		num++
		for k, v := range t.entities[i].tcpData {
			if _, ok := sum.tcpData[k]; !ok {
				sum.tcpData[k] = v
				continue
			}
			sum.tcpData[k] += v
		}
	}

	result := TCPCountDto{
		tcpData: map[string]int32{},
	}

	for k, v := range sum.tcpData {
		result.tcpData[k] = v / num
	}

	return result
}
