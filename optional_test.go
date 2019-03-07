package optional_test

import (
	op "github.com/MercuryThePlanet/optional"
	"testing"
)

const (
	TEST_STR = "test string"
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

func Test_Of(t *testing.T) {
	t.Run("Testing Of", Of_test)
	t.Run("Testing nil Of", OfNil_test)
}

func Of_test(t *testing.T) {
	defer shouldNotPanic("optional.Of with non nil value", t)

	var o *op.Optional
	o = op.Of(TEST_STR)

	if v, ok := o.Get().(string); !ok {
		t.Error("Returned type should be string.")
	} else if v != TEST_STR {
		t.Errorf("Expected `%v`, got `%v`", TEST_STR, v)
	}
}

func OfNil_test(t *testing.T) {
	defer shouldPanic("optional.Of with nil value", t)

	var _ *op.Optional = op.Of(nil)
	t.Fatal("this code should not be reachable.")
}

func Test_OfNilable(t *testing.T) {
	t.Run("Testing OfNilable", OfNilable_test)
	t.Run("Testing nil OfNilable", OfNilableNil_test)
}

func OfNilable_test(t *testing.T) {
	defer shouldNotPanic("optional.OfNilable", t)

	var o *op.Optional
	o = op.OfNilable(TEST_STR)

	if v, ok := o.Get().(string); !ok {
		t.Error("Returned type should be string.")
	} else if v != TEST_STR {
		t.Errorf("Expected `%v`, got `%v`", TEST_STR, v)
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

func Test_IfPresent(t *testing.T) {
	t.Run("Testing IfPresent", IfPresent_test)
	t.Run("Testing nil IfPresent not present", IfNotPresent_test)
}

func IfPresent_test(t *testing.T) {
	defer shouldNotPanic("optional.IfPresent", t)

	var o *op.Optional
	o = op.Of(TEST_STR)

	ok := false
	o.IfPresent(func(v op.T) {
		ok = true
	})

	if !ok {
		t.Error("IfPresent was not reached when it should have been.")
	}
}

func IfNotPresent_test(t *testing.T) {
	defer shouldNotPanic("optional.IfPresent", t)

	var o *op.Optional
	o = op.OfNilable(nil)

	o.IfPresent(func(v op.T) {
		t.Error("IfPresent was reached when it should not have been.")
	})
}
