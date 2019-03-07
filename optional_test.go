package optional_test

import (
	op "github.com/MercuryThePlanet/optional"
	"testing"
)

func shouldPanic(func_name string, t *testing.T) {
	if r := recover(); r == nil {
		t.Fatal(func_name + " should panic and did not.")
	}
}

func shouldNotPanic(func_name string, t *testing.T) {
	if r := recover(); r != nil {
		t.Fatalf(func_name+" should not panic and did: %v", r)
	}
}

func Test_OptionalOf(t *testing.T) {
	t.Run("Testing Of", Of_test)
	t.Run("Testing nil Of", OfNil_test)
}

func Of_test(t *testing.T) {
	defer shouldNotPanic("optional.Of with non nil value", t)

	test_val := "test"

	var o *op.Optional
	o = op.Of(test_val)

	if v, ok := o.Get().(string); !ok {
		t.Error("Returned type should be string.")
	} else if v != test_val {
		t.Errorf("Expected `%v`, got `%v`", test_val, v)
	}
}

func OfNil_test(t *testing.T) {
	defer shouldPanic("optional.Of with nil value", t)

	var _ *op.Optional = op.Of(nil)
	t.Fatal("this code should not be reachable.")
}

func Test_OptionalOfNilable(t *testing.T) {
	t.Run("Testing OfNilable", OfNullable_test)
	t.Run("Testing nil OfNilable", OfNullableNil_test)
}

func OfNilable_test(t *testing.T) {
	defer shouldNotPanic("optional.OfNilable", t)
	test_val := "test"

	var o *op.Optional
	o = op.OfNilable(test_val)

	if v, ok := o.Get().(string); !ok {
		t.Error("Returned type should be string.")
	} else if v != test_val {
		t.Errorf("Expected `%v`, got `%v`", test_val, v)
	}
}

func OfNilableNil_test(t *testing.T) {
	defer shouldNotPanic("optional.OfNilable", t)

	var o *op.Optional
	o = op.OfNilable(nil)

	if o.Get() != nil {
		t.Error("Returned type should be nil but is not.")
	}
}
