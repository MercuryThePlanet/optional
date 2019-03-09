[![GoDoc](https://godoc.org/github.com/MercuryThePlanet/optional?status.svg)](https://godoc.org/github.com/MercuryThePlanet/optional)
[![Coverage](./coverage_badge.png)](./coverage_badge.png)

# optional
package optional contains functions for creating an optional container.

An optional container may or may not contain a nil value. Optional methods allow you to write code without checking if a value is nil or not. If the value is nil the code will not be executed and enter a panic.

```go
package main

import (
	"fmt"
	op "github.com/MercuryThePlanet/optional"
)

func main() {
	op.OfNilable(nil).IfPresentOrElse(func(t op.T) {
		println("This code will not run.")
	}, func() {
		println("This code will run.")
	})

	for i := 0; i < 10; i++ {
		op.Of(i).Filter(func(t op.T) bool {
			return t.(int) > 5
		}).Map(func(t op.T) op.T {
			fmt.Printf("Original value is %d. ", t.(int))
			return t.(int) * 25
		}).IfPresent(func(t op.T) {
			fmt.Printf("Mapped value is %d.\n", t.(int))
		})
	}
}
```

Output:
```
This code will run.
Original value is 6. Mapped value is 150.
Original value is 7. Mapped value is 175.
Original value is 8. Mapped value is 200.
Original value is 9. Mapped value is 225.
```
