package memory

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Stats represents memory statistics for linux
//MemTotal: Total usable ram (i.e. physical ram minus a few reserved bits and the kernel binary code)
//MemFree: Is sum of LowFree+HighFree (overall stat)
//MemShared: 0; is here for compat reasons but always zero.
//Buffers: Memory in buffer cache. mostly useless as metric nowadays Relatively temporary storage for raw disk blocks shouldn’t get tremendously large (20MB or so)
//Cached: Memory in the pagecache (diskcache) minus SwapCache, Doesn’t include SwapCached
//SwapCache: Memory that once was swapped out, is swapped back in but still also is in the swapfile (if memory is needed it doesn’t need to be swapped out AGAIN because it is already in the swapfile. This saves I/O )

// Memory statistic
//HighTotal: is the total amount of memory in the high region. Highmem is all memory above (approx) 860MB of physical RAM. Kernel uses indirect tricks to access the high memory region. Data cache can go in this memory region.
//LowTotal: The total amount of non-highmem memory.
//LowFree: The amount of free memory of the low memory region. This is the memory the kernel can address directly. All kernel data structures need to go into low memory.
//SwapTotal: Total amount of physical swap memory.
//SwapFree: Total amount of swap memory free. Memory which has been evicted from RAM, and is temporarily on the disk
//Dirty: Memory which is waiting to get written back to the disk
//Writeback: Memory which is actively being written back to the disk
//Mapped: files which have been mapped, such as libraries
//Slab: in-kernel data structures cache
//Committed_AS: An estimate of how much RAM you would need to make a 99.99% guarantee that there never is OOM (out of memory) for this workload. Normally the kernel will overcommit memory. That means, say you do a 1GB malloc, nothing happens,really. Only when you start USING that malloc memory you will get real memory on demand, and just as much as you use. So you sort of take a mortgage and hope the bank doesn’t go bust. Other cases might include when you mmap a file that’s shared only when you write to it and you get a private copy of that data. While it normally is shared between processes. The Committed_AS is a guesstimate of how much RAM/swap you would need worst-case.
//PageTables: amount of memory dedicated to the lowest level of page tables.
//ReverseMaps: number of reverse mappings performed
//VmallocTotal: total size of vmalloc memory area
//VmallocUsed: amount of vmalloc area which is used
//VmallocChunk: largest contigious block of vmalloc area which is free
type Stats struct {
	Total, Used, Buffers, Cached, Free, Available, Active, Inactive,
	SwapTotal, SwapUsed, SwapCached, SwapFree, Mapped, Shmem, Slab,
	PageTables, Committed, VmallocUsed uint64
	MemAvailableEnabled bool
}

func Get() (*Stats, error) {
	file, err := os.Open("/proc/meminfo")
	defer file.Close()
	if err != nil {
		return nil, err
	}

	return collect(file)
}

func collect(out io.Reader) (*Stats, error) {
	scanner := bufio.NewScanner(out)
	var memory Stats
	memStats := map[string]*uint64{
		"MemTotal":     &memory.Total,
		"MemFree":      &memory.Free,
		"MemAvailable": &memory.Available,
		"Buffers":      &memory.Buffers,
		"Cached":       &memory.Cached,
		"Active":       &memory.Active,
		"Inactive":     &memory.Inactive,
		"SwapCached":   &memory.SwapCached,
		"SwapTotal":    &memory.SwapTotal,
		"SwapFree":     &memory.SwapFree,
		"Mapped":       &memory.Mapped,
		"Shmem":        &memory.Shmem,
		"Slab":         &memory.Slab,
		"PageTables":   &memory.PageTables,
		"Committed_AS": &memory.Committed,
		"VmallocUsed":  &memory.VmallocUsed,
	}
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i < 0 {
			continue
		}
		fld := line[:i]
		if ptr := memStats[fld]; ptr != nil {
			val := strings.TrimSpace(strings.TrimRight(line[i+1:], "kB"))
			if v, err := strconv.ParseUint(val, 10, 64); err == nil {
				*ptr = v * 1024
			}
			if fld == "MemAvailable" {
				memory.MemAvailableEnabled = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %s", err)
	}

	memory.SwapUsed = memory.SwapTotal - memory.SwapFree

	if memory.MemAvailableEnabled {
		memory.Used = memory.Total - memory.Available
	} else {
		memory.Used = memory.Total - memory.Free - memory.Buffers - memory.Cached
	}

	return &memory, nil
}
