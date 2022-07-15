# Image Optimizer
This package useful for metric

`
go get github.com/sonnam0904/metric
`


#Use
Example
```
package main

import (
	"fmt"
	"os"
    "github.com/sonnam0904/metric/cpu"
	"github.com/sonnam0904/metric/memory"
)

func main() {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	fmt.Printf("memory total: %d MB\n", memory.Total/1000000)
	fmt.Printf("memory used: %d MB\n", memory.Used/1000000)
	fmt.Printf("memory cached: %d MB\n", memory.Cached/1000000)
	fmt.Printf("memory free: %d MB\n", memory.Free/1000000)

    before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	total := float64(after.Total - before.Total)
	fmt.Printf("CPU user: %f %%\n", float64(after.User-before.User)/total*100)
	fmt.Printf("CPU system: %f %%\n", float64(after.System-before.System)/total*100)
	fmt.Printf("CPU idle: %f %%\n", float64(after.Idle-before.Idle)/total*100)
}
