package metric

import (
	"time"
	"github.com/sonnam0904/metric/cpu"
	"github.com/sonnam0904/metric/memory"
)

type Monitor struct {
	OsMemory OsMemory
	OsCPU OsCPU
}

type OsMemory struct {
	Total, 
	Used,
	Cached,
	Buffers,
	Free uint64
}

type OsCPU struct {
	Used,
	Idle,
	ProcsRunning,
	ProcsBlocked,
	System float64
	CPUCount int
}

func NewMonitor(duration int, info chan Monitor) {
	var m Monitor
	var interval = time.Duration(duration) * time.Second
	memory, err := memory.Get()
	cpu, err := cpu.Get()
	for {
		<-time.After(interval)
		// memory
		if err != nil {
			m.OsMemory.Total = uint64(0)
			m.OsMemory.Used = uint64(0)
			m.OsMemory.Cached = uint64(0)
			m.OsMemory.Buffers = uint64(0)
			m.OsMemory.Free = uint64(0)
		} else {
			m.OsMemory.Total = byteToMb(memory.Total)
			m.OsMemory.Used = byteToMb(memory.Used)
			m.OsMemory.Cached = byteToMb(memory.Cached)
			m.OsMemory.Buffers = byteToMb(memory.Buffers)
			m.OsMemory.Free = byteToMb(memory.Free)
		}

		// CPU
		if err != nil {
			m.OsCPU.Used = float64(0)
			m.OsCPU.System = float64(0)
			m.OsCPU.Idle = float64(0)
			m.OsCPU.CPUCount = int(0)
		} else {
			m.OsCPU.Used = float64(cpu.User)/float64(cpu.Total)*100
			m.OsCPU.Idle = float64(cpu.Idle)/float64(cpu.Total)*100
			m.OsCPU.System = float64(cpu.System)/float64(cpu.Total)*100
			m.OsCPU.ProcsRunning = float64(cpu.ProcsRunning)
			m.OsCPU.ProcsBlocked = float64(cpu.ProcsBlocked)
			m.OsCPU.CPUCount = int(cpu.CPUCount)
		}
		info<-m
		//b, _ := json.Marshal(m)
		//fmt.Println(string(b))
	}
}

func byteToMb(b uint64) uint64 {
	return b / 1024 / 1024
}