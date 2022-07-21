# Image Optimizer
This package useful for metric

`
go get github.com/sonnam0904/metric
`


#Use
Example
```
package metric

import (
	"encoding/json"
	"fmt"
	"github.com/sonnam0904/metric"
)

func main() {
	output := make(chan metric.Monitor)
	go metric.NewMonitor(2, output)
	for {
		select {
		case <-output:
			res, _ := json.Marshal(<-output)
			fmt.Println(string(res))
		}
	}
}
