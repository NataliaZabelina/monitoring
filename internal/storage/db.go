package storage

import "github.com/NataliaZabelina/monitoring/internal/storage/schema"

type Db struct {
	CPU_table    CPUTable
	System_table SystemTable
	Disk_table   DiskTable
	Socket_table SocketTable
	TCP_table    TCPTable
}

func (db *Db) Init() {
	db.CPU_table = (&schema.CpuLoadTable{}).Init()
	db.System_table = (&schema.SystemLoadTable{}).Init()
	db.Disk_table = (&schema.DiskLoadTable{}).Init()
	db.Socket_table = (&schema.SocketLoadTable{}).Init()
	db.TCP_table = (&schema.TCPCountTable{}).Init()
}

type CPUTable interface {
	AddEntity(schema.CpuLoadDto)
	GetAverage(period int32) schema.CpuLoadDto
}

type SystemTable interface {
	AddEntity(schema.SystemLoadDto)
	GetAverage(period int32) schema.SystemLoadDto
}

type DiskTable interface {
	AddEntity(schema.DiskLoadDto)
	GetAverage(period int32) schema.DiskLoadDto
}

type SocketTable interface {
	AddEntity(schema.SocketLoadInfo)
	GetAverage(period int32) schema.SocketLoadInfo
}

type TCPTable interface {
	AddEntity(schema.TCPCountDto)
	GetAverage(period int32) schema.TCPCountDto
}
