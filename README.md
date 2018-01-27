# Skinny
Lightweight in memory time-series database written in Golang

# Example

```go
package main

import (
	"fmt"

	"github.com/maurodelazeri/skinny"
)

func main() {
	testmetric := skinny.Metric{Capacity: 131487, Indexinterval: 1440} // 3 months of per minute, indexed on days
	testmetric.Init()

	test := map[string]interface{}{
		"test": "mauro",
	}

	testmetric.Insert(1416585010, test, true)
	testmetric.Insert(1416585011, test, true)
	testmetric.Insert(1416585012, test, true)
	testmetric.Insert(1416585012, test, true)
	testmetric.Insert(1416585012, test, true)
	testmetric.Insert(1416585013, test, true)
	testmetric.Insert(1416585014, test, true)
	testmetric.Insert(1416585015, test, true)
	testmetric.Insert(1416585015, test, true)

	for _, pnt := range testmetric.GetRange(1416585013, 1416585015) {
		fmt.Printf("Found: %+v\n", pnt)
	}
}

```

```
go run main.go 
Found: &{timestamp:1416585013 value:map[test:mauro] next:0xc42000a120}
Found: &{timestamp:1416585014 value:map[test:mauro] next:0xc42000a140}
Found: &{timestamp:1416585015 value:map[test:mauro] next:0xc42000a160}
Found: &{timestamp:1416585015 value:map[test:mauro] next:<nil>}
```
