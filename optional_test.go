package optional_test

import (
	op "github.com/MercuryThePlanet/optional"
	"testing"
)

func Test_OptionalOf(t *testing.T) {
	t.Run("Testing Of", OfTest)
	t.Run("Testing nil Of", OfNilTest)
}

func OfTest(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("optional.Of should not panic if value is not nil", r)
		}
	}()

	test_val := "test"
	of_op := op.Of(test_val)
	if v, ok := of_op.Get().(string); !ok {
		t.Error("Returned type should be string.")
	} else if v != test_val {
		t.Errorf("Expected `%v`, got `%v`", test_val, v)
	}
}
func OfNilTest(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("optional.Of should panic if nil is passed")
		}
	}()

	op.Of(nil)
}
