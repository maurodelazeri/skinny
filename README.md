# skinny
Lightweight in memory time-series database written in Golang

# Example

```go
package main

import (
	"fmt"
	"github.com/maurodelazeri/skinny"
)


func main() {
	testmetric := skinny.Metric{capacity: 131487, indexinterval: 1440} // 3 months of per minute, indexed on days
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
