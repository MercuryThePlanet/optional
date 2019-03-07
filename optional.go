package optional

type T = interface{}

type Optional struct {
	v   T
	set bool
}

type Interface interface {
	Equals(T) bool
}

type FlatMapFunc = func(T) *Optional
type MapFunc = func(T) T
type GetFunc = func() T

func Of(v T) *Optional {
	if v != nil {
		return &Optional{v: v, set: true}
	}
	panic("optional.Of takes a non-null value. Use OfNullable for potentially null values.")
}

func OfNullable(v T) *Optional {
	if v != nil {
		return &Optional{v: v, set: true}
	}
	return &Optional{set: false}
}

func (o *Optional) Get() T {
	if o.set {
		return o.v
	}
	return nil
}

func (o *Optional) IfPresent(f func(v T)) {
	if o.set {
		f(o.v)
	}
}

func (o *Optional) IsPresent() bool {
	return o.set
}

func (o *Optional) Map(f MapFunc) *Optional {
	if o.set {
		new_v := f(o.v)
		if new_v != nil {
			return &Optional{new_v, true}
		}
	}
	return &Optional{set: false}
}

func (o *Optional) FlatMap(f FlatMapFunc) *Optional {
	if o.set {
		new_v := f(o.v)
		if new_v != nil {
			return new_v
		}
	}
	return &Optional{set: false}
}

func (o *Optional) OrElse(other T) T {
	if o.set {
		return o.v
	} else {
		return other
	}
}

func (o *Optional) OrElseGet(f GetFunc) T {
	if o.set {
		return o.v
	} else {
		return f()
	}
}

func (o *Optional) OrElsePanic(p string) T {
	if o.set {
		return o.v
	} else {
		panic(p)
	}
}
