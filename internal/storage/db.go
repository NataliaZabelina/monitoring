package storage

import "github.com/NataliaZabelina/monitoring/internal/storage/schema"

type DB struct {
	CPUTable    CPUTable
	SystemTable SystemTable
	DiskTable   DiskTable
	// SocketTable SocketTable
	// TCPTable    TCPTable
}

func (db *DB) Init() {
	db.CPUTable = (&schema.CPULoadTable{}).Init()
	db.SystemTable = (&schema.SystemLoadTable{}).Init()
	db.DiskTable = (&schema.DiskTable{}).Init()
	// db.SocketTable = (&schema.SocketLoadTable{}).Init()
	// db.TCPTable = (&schema.TCPCountTable{}).Init()
}

type CPUTable interface {
	AddEntity(schema.CPULoadDto)
	GetAverage(period int32) schema.CPULoadDto
}

type SystemTable interface {
	AddEntity(schema.SystemLoadDto)
	GetAverage(period int32) schema.SystemLoadDto
}

type DiskTable interface {
	AddEntity(schema.DiskDto)
	GetAverage(period int32) schema.DiskDto
}

// type DiskFSTable interface {
// 	AddEntity(schema.DiskFSDto)
// 	GetAverage(period int32) schema.DiskFSDto
// }

// type SocketTable interface {
// 	AddEntity(schema.SocketLoadInfo)
// 	GetAverage(period int32) schema.SocketLoadInfo
// }

type TCPTable interface {
	AddEntity(schema.TCPCountDto)
	GetAverage(period int32) schema.TCPCountDto
}
