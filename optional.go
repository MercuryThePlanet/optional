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

type FlatMapFunc = func(T) *Optional
type FilterFunc = func(T) bool
type MapFunc = func(T) T
type GetFunc = func() T
type Supplier = func(Ts) *Optional

func (o *Optional) Or(f Supplier, ts ...T) *Optional {
	if o.present {
		return o
	} else {
		return f(ts)
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

func (o *Optional) Filter(f FilterFunc) *Optional {
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

func (o *Optional) set(t interface{}, present bool) *Optional {
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

func (o *Optional) IfPresent(f func(t T)) {
	if o.present {
		f(o.t)
	}
}

func (o *Optional) IfPresentOrElse(f func(t T), other func()) {
	if o.present {
		f(o.t)
	} else {
		other()
	}
}

func (o *Optional) IsPresent() bool {
	return o.present
}

func (o *Optional) Map(f MapFunc) *Optional {
	if o.present {
		new_t := f(o.t)
		return o.set(new_t, new_t != nil)
	}
	return o.set(nil, false)
}

func (o *Optional) FlatMap(f FlatMapFunc) *Optional {
	if o.present {
		new_t := f(o.t)
		if new_t != nil {
			return new_t
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

func (o *Optional) OrElseGet(f GetFunc) T {
	if o.present {
		return o.t
	} else {
		return f()
	}
}

func (o *Optional) OrElsePanic(p string) T {
	if o.present {
		return o.t
	}
	panic(p)
}
