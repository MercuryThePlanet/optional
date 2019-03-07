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

type S struct{ v int }

func (s *S) Cmpr(t op.T) int {
	o, ok := t.(*S)
	if !ok {
		return -2
	}
	if s.v > o.v {
		return 1
	} else if s.v < o.v {
		return -1
	} else {
		return 0
	}
}

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

func Test_Empty(t *testing.T) {
	t.Run("Empty", Empty_test)
}

func Empty_test(t *testing.T) {
	defer shouldNotPanic("optional.Empty", t)

	if op.Empty().IsPresent() {
		t.Error("Empty optional should have no value")
	}
}

func Test_Equals(t *testing.T) {
	t.Run("Equals", Equals_test)
	t.Run("Equals not equal", EqualsNot_test)
	t.Run("Equals Cmpr not implemented", EqualsNoCmpr_test)
	t.Run("Equals other not Optional", EqualsNotOptional_test)
	t.Run("Equals value not present", EqualsNotPresent_test)
	t.Run("Equals other value not present", EqualsOtherNotPresent_test)
}

func Equals_test(t *testing.T) {
	defer shouldNotPanic("optional.Equals", t)

	o := op.Of(&S{})
	other := op.Of(&S{})
	if !o.Equals(other) {
		t.Error("The values passed should be equal.")
	}
}

func EqualsNot_test(t *testing.T) {
	defer shouldNotPanic("optional.Equals", t)

	o := op.Of(&S{1})
	other := op.Of(&S{-1})
	if o.Equals(other) {
		t.Error("The values passed should not be equal.")
	}
}

func EqualsNoCmpr_test(t *testing.T) {
	defer shouldNotPanic("optional.Equals", t)

	o := op.Of(struct{}{})
	other := op.Of(&S{-1})
	if o.Equals(other) {
		t.Error("The values passed should not be equal.")
	}
}

func EqualsNotOptional_test(t *testing.T) {
	defer shouldNotPanic("optional.Equals", t)

	o := op.Of(&S{})
	if o.Equals(&S{}) {
		t.Error("The values passed should not be equal.")
	}
}

func EqualsNotPresent_test(t *testing.T) {
	defer shouldNotPanic("optional.Equals", t)

	if op.Empty().Equals(&S{}) {
		t.Error("The values passed should not be equal.")
	}
}

func EqualsOtherNotPresent_test(t *testing.T) {
	defer shouldNotPanic("optional.Equals", t)

	if op.Of(&S{}).Equals(op.Empty()) {
		t.Error("The values passed should not be equal.")
	}
}

func Test_Filter(t *testing.T) {
	t.Run("Filter", Filter_test)
	t.Run("Filter remove", FilterRemove_test)
}

func Filter_test(t *testing.T) {
	defer shouldNotPanic("optional.Filter", t)

	var o *op.Optional
	o = op.OfNilable(TEST_STR).Filter(func(v op.T) bool {
		_, ok := v.(string)
		return ok
	})
	if !o.IsPresent() {
		t.Error("Value should be present.")
	} else if v := o.Get().(string); v != TEST_STR {
		t.Errorf("Expected `%v`, got `%v`", TEST_STR, v)
	}
}

func FilterRemove_test(t *testing.T) {
	defer shouldNotPanic("optional.Filter", t)

	var o *op.Optional
	o = op.OfNilable(TEST_STR).Filter(func(v op.T) bool {
		_, ok := v.(int)
		return ok
	})
	if o.IsPresent() {
		t.Error("Value should not be present.")
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

func Test_Or(t *testing.T) {
	t.Run("Or", Or_test)
	t.Run("Or Other", OrOther_test)
	t.Run("Pass supplier single param", OrOtherSingleParam_test)
	t.Run("Pass supplier multiple params", OrOtherMultipleParams_test)
}

func Or_test(t *testing.T) {
	defer shouldNotPanic("optional.Or", t)

	var o *op.Optional
	o = op.OfNilable(TEST_INT).Or(func(ts op.Ts) *op.Optional {
		return op.Of(TEST_STR)
	})

	if v := o.Get().(int); v != TEST_INT {
		t.Errorf("Expected `%v`, got `%v`", TEST_INT, v)
	}
}

func OrOther_test(t *testing.T) {
	defer shouldNotPanic("optional.Or", t)

	var o *op.Optional
	o = op.OfNilable(nil).Or(func(ts op.Ts) *op.Optional {
		return op.Of(TEST_STR)
	})

	if v := o.Get().(string); v != TEST_STR {
		t.Errorf("Expected `%v`, got `%v`", TEST_STR, v)
	}
}

func OrOtherSingleParam_test(t *testing.T) {

	var o *op.Optional
	o = op.OfNilable(nil).Or(func(ts op.Ts) *op.Optional {
		return op.Of(ts[0].(string))
	}, TEST_STR)

	if v := o.Get().(string); v != TEST_STR {
		t.Errorf("Expected `%v`, got `%v`", TEST_STR, v)
	}
}

func OrOtherMultipleParams_test(t *testing.T) {

	var o *op.Optional
	o = op.OfNilable(nil).Or(func(ts op.Ts) *op.Optional {
		return op.Of(ts[0].(int) + ts[1].(int))
	}, TEST_INT, TEST_OTHER)

	expected := TEST_INT + TEST_OTHER
	if v := o.Get().(int); v != expected {
		t.Errorf("Expected `%v`, got `%v`", expected, v)
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

func Test_IfPresentOrElse(t *testing.T) {
	t.Run("IfPresentOrElse", IfPresentOrElse_test)
	t.Run("IfPresentOrElse not present", IfPresentOrElseOther_test)
}

func IfPresentOrElse_test(t *testing.T) {
	defer shouldNotPanic("optional.IfPresentOrElse", t)

	var o *op.Optional
	o = op.Of(TEST_STR)

	o.IfPresentOrElse(func(v op.T) {
	}, func() {
		t.Error("IfPresentOrElse other was reached when it should not have been.")
	})
}

func IfPresentOrElseOther_test(t *testing.T) {
	defer shouldNotPanic("optional.IfPresentOrElse", t)

	var o *op.Optional
	o = op.OfNilable(nil)

	o.IfPresentOrElse(func(v op.T) {
		t.Error("IfPresentOrElse first was reached when it should not have been.")
	}, func() {
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
