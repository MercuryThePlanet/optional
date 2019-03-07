// package optional contains functions for creating an optional container.
//
// An optional container may or may not contain a nil value. Optional methods
// allow you to write code without checking if a value is nil or not. If the
// value is nil the code will not be executed and enter a panic.
package optional

// Alias of empty interface for ease of use
type T = interface{}

// Alias of slice of empty interface for ease of use
type Ts = []T

// struct Optional is the container struct.
type Optional struct {
	t       T
	present bool
}

// Interface contains one method, Cmpr, which takes an empty interface
// and returns an integer. If the integer is 0, the passed value is
// considered equal to the caller. This method must be implemented for
// the Optional.Equals() method to work.
type Interface interface {
	Cmpr(T) int
}

// A Predicate function signature.
//
// Takes an empty interface, runs some logic and returns a bool.
type Predicate = func(T) bool

// A Supplier function signature.
//
// A variadic function that returns an empty interface.
type Supplier = func(Ts) T

// A Mapper function signature.
type Mapper = func(T) T

// A Consumer function signature.
type Consumer = func(T)

// A Runnable function signature.
type Runnable = func()

// Returns an empty Optional instance.
func Empty() *Optional {
	return &Optional{}
}

func (o *Optional) set(t T, present bool) *Optional {
	o.t = t
	o.present = present
	return o
}

// Indicates if another object is equal to this Optional.
//
// Two objects are considered equal if:
// - They are both optionals.
// - They are both empty optionals or;
// - The contained values return zero after running the Cmpr method.
func (o *Optional) Equals(t T) bool {
	if o.IsPresent() {
		i, isInterface := o.t.(Interface)
		other, ok := t.(*Optional)
		if ok && isInterface && other.IsPresent() {
			return i.Cmpr(other.Get()) == 0
		}
	}
	return false
}

// If a value is present, returns the result of applying the given
// Optional-bearing mapping function to the value, otherwise returns an empty
// Optional.
func (o *Optional) Filter(f Predicate) *Optional {
	if f(o.t) {
		return o
	} else {
		return o.set(nil, false)
	}
}

// Returns an Optional describing the given non-nil value.
func Of(t T) *Optional {
	if t != nil {
		return &Optional{t: t, present: true}
	}
	panic("optional.Of takes a non-nil talue. Use OfNilable for potentially nil talues.")
}

// Returns an Optional describing the given value, if non-nil, otherwise
// returns an empty Optional.
func OfNilable(t T) *Optional {
	return &Optional{t: t, present: t != nil}
}

func OfErrorable(t T, err error) *Optional {
	if err == nil {
		return &Optional{t: t, present: t != nil}
	}
	return &Optional{}
}

// If a value is present, returns an Optional describing the value, otherwise
// returns an Optional produced by the supplying function.
func (o *Optional) Or(f Supplier, ts ...T) *Optional {
	if o.present {
		return o
	} else {
		t := f(ts)
		return o.set(t, t != nil)
	}
}

// If a value is present, returns the value, otherwise returns nil.
func (o *Optional) Get() T {
	if o.present {
		return o.t
	}
	return nil
}

// If a value is present, performs the given action with the value, otherwise
// does nothing.
func (o *Optional) IfPresent(f Consumer) {
	if o.present {
		f(o.t)
	}
}

// If a value is present, performs the given action with the value, otherwise
// performs the given runnable action.
func (o *Optional) IfPresentOrElse(f Consumer, other Runnable) {
	if o.present {
		f(o.t)
	} else {
		other()
	}
}

// If a value is present, returns true, otherwise false.
func (o *Optional) IsPresent() bool {
	return o.present
}

// If a value is present, returns an Optional describing (as if by ofNullable(T))
// the result of applying the given mapping function to the value, otherwise
// returns an empty Optional.
func (o *Optional) Map(f Mapper) *Optional {
	if o.present {
		mapped_t := f(o.t)
		return o.set(mapped_t, mapped_t != nil)
	}
	return o.set(nil, false)
}

// If a value is present, returns the result of applying the given
// Optional-bearing mapping function to the value, otherwise returns an empty
// Optional.
func (o *Optional) FlatMap(f Mapper) T {
	if o.present {
		mapped_t := f(o.t)
		if mapped_t != nil {
			return mapped_t
		}
	}
	return o.set(nil, false)
}

// If a value is present, returns the value, otherwise returns other.
func (o *Optional) OrElse(other T) T {
	if o.present {
		return o.t
	} else {
		return other
	}
}

// If a value is present, returns the value, otherwise returns the result
// produced by the supplying function.
func (o *Optional) OrElseGet(f Supplier, ts ...T) T {
	if o.present {
		return o.t
	} else {
		return f(ts)
	}
}

// If a value is present, returns the value, otherwise panics.
func (o *Optional) OrElsePanic(p string) T {
	if o.present {
		return o.t
	}
	panic(p)
}
