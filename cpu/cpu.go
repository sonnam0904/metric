package cpu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Stats represents cpu statistics for linux
type Stats struct {
	User, Nice, System, Idle, Iowait, Irq, Softirq, Steal, Guest, GuestNice, Total, ProcsRunning, ProcsBlocked uint64
	CPUCount, StatCount                                                            int
}

type cpuStat struct {
	name string
	ptr  *uint64
}

func Get() (*Stats, error) {
	file, err := os.Open("/proc/stat")
	defer file.Close()
	if err != nil {
		return nil, err
	}
	return collect(file)
}

func collect(out io.Reader) (*Stats, error) {
	scanner := bufio.NewScanner(out)
	var cpu Stats

	//user: normal processes executing in user mode
	//nice: niced processes executing in user mode
	//system: processes executing in kernel mode
	//idle: twiddling thumbs
	//iowait: waiting for I/O to complete
	//irq: servicing interrupts
	//softirq: servicing softirqs

	cpuStats := []cpuStat{
		{"user", &cpu.User},
		{"nice", &cpu.Nice},
		{"system", &cpu.System},
		{"idle", &cpu.Idle},
		{"iowait", &cpu.Iowait},
		{"irq", &cpu.Irq},
		{"softirq", &cpu.Softirq},
		{"steal", &cpu.Steal},
		{"guest", &cpu.Guest},
		{"guest_nice", &cpu.GuestNice},
		{"procs_running", &cpu.ProcsRunning},
		{"procs_blocked", &cpu.ProcsBlocked},
	}

	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to scan")
	}

	val := strings.Fields(scanner.Text())[1:]

	cpu.StatCount = len(val)
	for i, valStr := range val {

		val, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to scan %s ", cpuStats[i].name)
		}
		*cpuStats[i].ptr = val
		cpu.Total += val
	}

	// https://github.com/torvalds/linux/blob/4ec9f7a18/kernel/sched/cputime.c#L151-L158
	cpu.Total -= cpu.Guest
	cpu.Total -= cpu.GuestNice
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ' ')
		if i < 0 {
			continue
		}
		k := line[:i]
		v := line[i+1:]
		if strings.HasPrefix(k, "cpu") {
			cpu.CPUCount++
		}

		switch k {
			case "procs_running":
			if vi, err := strconv.ParseUint(v, 10, 64); err == nil {
				cpu.ProcsRunning = vi
			}
			case "procs_blocked":
			if vi, err := strconv.ParseUint(v, 10, 64); err == nil {
				cpu.ProcsBlocked = vi
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %s", err)
	}

	return &cpu, nil
}