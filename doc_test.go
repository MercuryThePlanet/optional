package optional_test

import (
	"fmt"
	op "github.com/MercuryThePlanet/optional"
	"strconv"
)

type C struct {
	str string
}

type B struct {
	c *C
}

type A struct {
	b *B
}

func (c *C) Cmpr(t op.T) int {
	if v, ok := t.(string); !ok {
		return -2
	} else if v == c.str {
		return 0
	} else {
		return -1
	}
}

func Example() {
	str_value := "00001234.43210000"

	a := &A{&B{&C{str_value}}}

	cOp := op.OfNilable(a).Map(func(t op.T) op.T {
		return t.(*A).b
	}).Map(func(t op.T) op.T {
		return t.(*B).c
	})

	fmt.Printf("Values are equal: %t\n", cOp.Equals(op.Of(str_value)))

	str := cOp.FlatMap(func(t op.T) op.T {
		return t.(*C).str
	}).(string)

	floatOp := op.OfErrorable(strconv.ParseFloat(str, 64))

	floatOp.IfPresent(func(t op.T) {
		fmt.Printf("Value is: %f\n", t.(float64))
	})

	op.Of(floatOp.Get().(float64)).Filter(func(t op.T) bool {
		return t.(float64) < 1234.0
	}).IfPresentOrElse(func(t op.T) {
		fmt.Println("This code will not be reached.")
	}, func() {
		fmt.Println("In runnable.")
	})

	fmt.Printf("%f is the same as %f\n", floatOp.Get().(float64),
		floatOp.OrElse(4321.001).(float64))
	fmt.Printf("but %f is not the same as %f\n", floatOp.Get().(float64),
		op.Empty().OrElse(4321.001).(float64))

	fmt.Printf("Again, %f is the same as %f\n", floatOp.Get().(float64),
		floatOp.OrElseGet(func(_ op.Ts) op.T {
			return 4321.001
		}).(float64))
	fmt.Printf("but %f is not the same as %f\n", floatOp.Get().(float64),
		op.Empty().OrElseGet(func(_ op.Ts) op.T {
			return 4321.001
		}).(float64))

	int_val := op.OfNilable(nil).Or(func(ts op.Ts) op.T {
		return ts[0].(int) + ts[1].(int)
	}, 12345, 54321).OrElsePanic("This will not panic").(int)

	fmt.Printf("The or int value is %d\n", int_val)
}
