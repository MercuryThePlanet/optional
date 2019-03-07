package optional

type T = interface{}

type Optional struct {
	v       T
	present bool
}

type Interface interface {
	Equals(T) bool
}

type FlatMapFunc = func(T) *Optional
type FilterFunc = func(T) bool
type MapFunc = func(T) T
type GetFunc = func() T

func Empty() *Optional {
	return &Optional{}
}

func (o *Optional) Filter(f FilterFunc) *Optional {
	if f(o.v) {
		return o
	} else {
		return o.set(nil, false)
	}
}

func Of(v T) *Optional {
	if v != nil {
		return &Optional{v: v, present: true}
	}
	panic("optional.Of takes a non-nil value. Use OfNilable for potentially nil values.")
}

func OfNilable(v T) *Optional {
	return &Optional{v: v, present: v != nil}
}

func (o *Optional) set(v interface{}, present bool) *Optional {
	o.v = v
	o.present = present
	return o
}

func (o *Optional) Get() T {
	if o.present {
		return o.v
	}
	return nil
}

func (o *Optional) IfPresent(f func(v T)) {
	if o.present {
		f(o.v)
	}
}

func (o *Optional) IfPresentOrElse(f func(v T), other func()) {
	if o.present {
		f(o.v)
	} else {
		other()
	}
}

func (o *Optional) IsPresent() bool {
	return o.present
}

func (o *Optional) Map(f MapFunc) *Optional {
	if o.present {
		new_v := f(o.v)
		return o.set(new_v, new_v != nil)
	}
	return o.set(nil, false)
}

func (o *Optional) FlatMap(f FlatMapFunc) *Optional {
	if o.present {
		new_v := f(o.v)
		if new_v != nil {
			return new_v
		}
	}
	return o.set(nil, false)
}

func (o *Optional) OrElse(other T) T {
	if o.present {
		return o.v
	} else {
		return other
	}
}

func (o *Optional) OrElseGet(f GetFunc) T {
	if o.present {
		return o.v
	} else {
		return f()
	}
}

func (o *Optional) OrElsePanic(p string) T {
	if o.present {
		return o.v
	}
	panic(p)
}
