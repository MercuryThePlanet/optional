// package optional contains functions for creating an optional container.
//
// An optional container may or may not contain a nil value. Optional methods
// allow you to write code without checking if a value is nil or not. If the
// value is nil the code will not be executed and enter a panic.
package optional

type T = interface{}
type Ts = []T

type Optional struct {
	t       T
	present bool
}

type Interface interface {
	Cmpr(T) int
}

type Predicate = func(T) bool
type Supplier = func(Ts) T
type Mapper = func(T) T
type Consumer = func(T)
type Runnable = func()

func (o *Optional) Or(f Supplier, ts ...T) *Optional {
	if o.present {
		return o
	} else {
		t := f(ts)
		return o.set(t, t != nil)
	}
}

func Empty() *Optional {
	return &Optional{}
}

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

func (o *Optional) Filter(f Predicate) *Optional {
	if f(o.t) {
		return o
	} else {
		return o.set(nil, false)
	}
}

func Of(t T) *Optional {
	if t != nil {
		return &Optional{t: t, present: true}
	}
	panic("optional.Of takes a non-nil talue. Use OfNilable for potentially nil talues.")
}

func OfNilable(t T) *Optional {
	return &Optional{t: t, present: t != nil}
}

func (o *Optional) set(t T, present bool) *Optional {
	o.t = t
	o.present = present
	return o
}

func (o *Optional) Get() T {
	if o.present {
		return o.t
	}
	return nil
}

func (o *Optional) IfPresent(f Consumer) {
	if o.present {
		f(o.t)
	}
}

func (o *Optional) IfPresentOrElse(f Consumer, other Runnable) {
	if o.present {
		f(o.t)
	} else {
		other()
	}
}

func (o *Optional) IsPresent() bool {
	return o.present
}

func (o *Optional) Map(f Mapper) *Optional {
	if o.present {
		mapped_t := f(o.t)
		return o.set(mapped_t, mapped_t != nil)
	}
	return o.set(nil, false)
}

func (o *Optional) FlatMap(f Mapper) T {
	if o.present {
		mapped_t := f(o.t)
		if mapped_t != nil {
			return mapped_t
		}
	}
	return o.set(nil, false)
}

func (o *Optional) OrElse(other T) T {
	if o.present {
		return o.t
	} else {
		return other
	}
}

func (o *Optional) OrElseGet(f Supplier, ts ...T) T {
	if o.present {
		return o.t
	} else {
		return f(ts)
	}
}

func (o *Optional) OrElsePanic(p string) T {
	if o.present {
		return o.t
	}
	panic(p)
}
