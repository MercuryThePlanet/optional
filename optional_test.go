package optional_test

import (
	op "github.com/MercuryThePlanet/optional"
	"strconv"
	"testing"
)

const (
	TEST_STR   = "123"
	TEST_INT   = 123
	TEST_OTHER = 321
	TEST_PANIC = "panicking"
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
	t.Run("Of", Of_test)
	t.Run("nil Of", OfNil_test)
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
	t.Run("OfNilable", OfNilable_test)
	t.Run("nil OfNilable", OfNilableNil_test)
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
	t.Run("IfPresent", IfPresent_test)
	t.Run("IfPresent not present", IfNotPresent_test)
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

func Test_IsPresent(t *testing.T) {
	t.Run("IsPresent", IsPresent_test)
	t.Run("IsPresent not present", IsNotPresent_test)
}

func IsPresent_test(t *testing.T) {
	defer shouldNotPanic("optional.IsPresent", t)

	var o *op.Optional
	o = op.Of(TEST_STR)

	if !o.IsPresent() {
		t.Error("IsPresent was false when it should have true.")
	}
}

func IsNotPresent_test(t *testing.T) {
	defer shouldNotPanic("optional.IsPresent", t)

	var o *op.Optional
	o = op.OfNilable(nil)

	if o.IsPresent() {
		t.Error("IsPresent was true when it should have false.")
	}
}

func Test_Map(t *testing.T) {
	t.Run("Map", Map_test)
	t.Run("Map chain", MapChain_test)
	t.Run("Map returns nil", MapReturnsNil_test)
	t.Run("Map empty optional", MapEmptyOptional_test)
}

func Map_test(t *testing.T) {
	defer shouldNotPanic("optional.Map", t)

	var o *op.Optional
	o = op.Of(TEST_STR).Map(func(v op.T) op.T {
		val, err := strconv.Atoi(v.(string))
		if err != nil {
			t.Fatalf("Could not map value from string to int.")
		}
		return val
	})

	if v, ok := o.Get().(int); !ok {
		t.Error("Returned type should be int.")
	} else if v != TEST_INT {
		t.Errorf("Expected `%v`, got `%v`", TEST_INT, v)
	}
}

func MapChain_test(t *testing.T) {
	defer shouldNotPanic("optional.Map", t)

	var o *op.Optional
	o = op.Of(TEST_STR).Map(func(v op.T) op.T {
		val, err := strconv.Atoi(v.(string))
		if err != nil {
			t.Fatalf("Could not map value from string to int.")
		}
		return val
	}).Map(func(v op.T) op.T {
		return strconv.Itoa(v.(int))
	})

	if v, ok := o.Get().(string); !ok {
		t.Error("Returned type should be string.")
	} else if v != TEST_STR {
		t.Errorf("Expected `%v`, got `%v`", TEST_STR, v)
	}
}

func MapReturnsNil_test(t *testing.T) {
	defer shouldNotPanic("optional.Map", t)

	var o *op.Optional
	o = op.Of(TEST_STR).Map(func(v op.T) op.T {
		return nil
	})

	if o.Get() != nil {
		t.Error("Optional should be empty.")
	}
}

func MapEmptyOptional_test(t *testing.T) {
	defer shouldNotPanic("optional.Map", t)

	var o *op.Optional
	o = op.OfNilable(nil)

	o.Map(func(v op.T) op.T {
		t.Fatal("Map on empty optional should not run")
		return nil
	})
}

func Test_FlatMap(t *testing.T) {
	t.Run("FlatMap", FlatMap_test)
	t.Run("FlatMap chain", FlatMapChain_test)
	t.Run("FlatMap returns nil", FlatMapReturnsNil_test)
	t.Run("Map empty optional", FlatMapEmptyOptional_test)
}

func FlatMap_test(t *testing.T) {
	defer shouldNotPanic("optional.FlatMap", t)

	var o *op.Optional
	o = op.Of(TEST_INT).FlatMap(func(v op.T) *op.Optional {
		return op.Of(v)
	})

	if v, ok := o.Get().(int); !ok {
		t.Error("Returned type should be int.")
	} else if v != TEST_INT {
		t.Errorf("Expected `%v`, got `%v`", TEST_INT, v)
	}
}

func FlatMapChain_test(t *testing.T) {
	defer shouldNotPanic("optional.FlatMap", t)

	var o *op.Optional
	o = op.Of(TEST_INT).FlatMap(func(v op.T) *op.Optional {
		return op.Of(strconv.Itoa(v.(int)))
	}).FlatMap(func(v op.T) *op.Optional {
		val, err := strconv.Atoi(v.(string))
		if err != nil {
			t.Fatalf("Could not map value from string to int.")
		}
		return op.Of(val)
	})

	if v, ok := o.Get().(int); !ok {
		t.Error("Returned type should be int.")
	} else if v != TEST_INT {
		t.Errorf("Expected `%v`, got `%v`", TEST_INT, v)
	}
}

func FlatMapReturnsNil_test(t *testing.T) {
	defer shouldNotPanic("optional.FlatMap", t)

	var o *op.Optional
	o = op.Of(TEST_STR).FlatMap(func(v op.T) *op.Optional {
		return nil
	})

	if o.Get() != nil {
		t.Error("Optional should be empty.")
	}
}

func FlatMapEmptyOptional_test(t *testing.T) {
	defer shouldNotPanic("optional.FlatMap", t)

	var o *op.Optional
	o = op.OfNilable(nil).FlatMap(func(v op.T) *op.Optional {
		t.Fatal("FlatMap on empty optional should not run")
		return nil
	})

	if o.Get() != nil {
		t.Error("Optional should be empty.")
	}
}

func Test_OrElse(t *testing.T) {
	t.Run("OrElse", OrElse_test)
	t.Run("OrElse Other", OrElseOther_test)
	t.Run("OrElse Nil Other", OrElseNilOther_test)
	t.Run("OrElse Different Types", OrElseDifOther_test)
}

func OrElse_test(t *testing.T) {
	defer shouldNotPanic("optional.OrElse", t)

	var o *op.Optional
	o = op.Of(TEST_INT)

	if v, ok := o.OrElse(TEST_OTHER).(int); !ok {
		t.Error("Returned type should be int.")
	} else if v != TEST_INT {
		t.Errorf("Expected `%v`, got `%v`", TEST_INT, v)
	}
}

func OrElseOther_test(t *testing.T) {
	defer shouldNotPanic("optional.OrElse", t)

	var o *op.Optional
	o = op.OfNilable(nil)

	if v, ok := o.OrElse(TEST_OTHER).(int); !ok {
		t.Error("Returned type should be int.")
	} else if v != TEST_OTHER {
		t.Errorf("Expected `%v`, got `%v`", TEST_OTHER, v)
	}
}

func OrElseNilOther_test(t *testing.T) {
	defer shouldNotPanic("optional.OrElse", t)

	var o *op.Optional
	o = op.OfNilable(nil)

	if o.OrElse(nil) != nil {
		t.Error("Returned type should be nil.")
	}
}

func OrElseDifOther_test(t *testing.T) {
	defer shouldNotPanic("optional.OrElse", t)

	switch v := op.Of(TEST_STR).OrElse(TEST_INT).(type) {
	case string:
		if v != TEST_STR {
			t.Errorf("Expected `%v`, got `%v`", TEST_STR, v)
		}
	case int:
		t.Error("Returned other: type should be string.")
	default:
		t.Error("Returned type should be string.")
	}

	switch v := op.OfNilable(nil).OrElse(TEST_INT).(type) {
	case string:
		t.Error("Returned type should be int.")
	case int:
		if v != TEST_INT {
			t.Errorf("Expected `%v`, got `%v`", TEST_INT, v)
		}
	default:
		t.Error("Returned type should be int.")
	}
}

func Test_OrElseGet(t *testing.T) {
	t.Run("OrElseGet", OrElseGet_test)
	t.Run("OrElseGet Other", OrElseGetOther_test)
	t.Run("OrElseGet Nil Other", OrElseGetNilOther_test)
	t.Run("OrElseGet Different Types", OrElseGetDifOther_test)
}

func OtherFunc() op.T {
	return TEST_OTHER
}

func OrElseGet_test(t *testing.T) {
	defer shouldNotPanic("optional.OrElseGet", t)

	var o *op.Optional
	o = op.Of(TEST_INT)

	if v, ok := o.OrElseGet(OtherFunc).(int); !ok {
		t.Error("Returned type should be int.")
	} else if v != TEST_INT {
		t.Errorf("Expected `%v`, got `%v`", TEST_INT, v)
	}
}

func OrElseGetOther_test(t *testing.T) {
	defer shouldNotPanic("optional.OrElseGet", t)

	var o *op.Optional
	o = op.OfNilable(nil)

	if v, ok := o.OrElseGet(OtherFunc).(int); !ok {
		t.Error("Returned type should be int.")
	} else if v != TEST_OTHER {
		t.Errorf("Expected `%v`, got `%v`", TEST_OTHER, v)
	}
}

func OrElseGetNilOther_test(t *testing.T) {
	defer shouldNotPanic("optional.OrElseGet", t)

	var o *op.Optional
	o = op.OfNilable(nil)

	if o.OrElseGet(func() op.T { return nil }) != nil {
		t.Error("Returned type should be nil.")
	}
}

func OrElseGetDifOther_test(t *testing.T) {
	defer shouldNotPanic("optional.OrElseGet", t)

	switch v := op.Of(TEST_STR).OrElseGet(OtherFunc).(type) {
	case string:
		if v != TEST_STR {
			t.Errorf("Expected `%v`, got `%v`", TEST_STR, v)
		}
	case int:
		t.Error("Returned other: type should be string.")
	default:
		t.Error("Returned type should be string.")
	}

	switch v := op.OfNilable(nil).OrElseGet(OtherFunc).(type) {
	case string:
		t.Error("Returned type should be int.")
	case int:
		if v != TEST_OTHER {
			t.Errorf("Expected `%v`, got `%v`", TEST_OTHER, v)
		}
	default:
		t.Error("Returned type should be int.")
	}
}

func Test_OrElsePanic(t *testing.T) {
	t.Run("OrElsePanic", OrElsePanic_test)
	t.Run("OrElsePanic panic", OrElsePanicOther_test)
}

func OrElsePanic_test(t *testing.T) {
	defer shouldNotPanic("optional.OrElsePanic", t)

	var o *op.Optional
	o = op.Of(TEST_INT)

	if v, ok := o.OrElsePanic(TEST_PANIC).(int); !ok {
		t.Error("Returned type should be int.")
	} else if v != TEST_INT {
		t.Errorf("Expected `%v`, got `%v`", TEST_INT, v)
	}
}

func OrElsePanicOther_test(t *testing.T) {
	defer shouldPanic("optional.OrElsePanic", t)

	op.OfNilable(nil).OrElsePanic(TEST_PANIC)
	t.Fatal("This code should be unreachable.")
}
