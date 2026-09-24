package utils

type ContainChecker[T comparable, A any] struct {
	dict map[T]A
	fits bool
}

func NewContainChecker[T comparable, A any](dict map[T]A) *ContainChecker[T, A] {
	return &ContainChecker[T, A]{
		dict: dict,
		fits: true,
	}
}

func (a *ContainChecker[T, A]) Any(objs... T) *ContainChecker[T, A] {
	if !a.fits { return a }

	for _, o := range objs {
		if _, ok := a.dict[o]; ok {
			return a
		}
	}

	a.fits = false
	return a
}

func (a *ContainChecker[T, A]) All(objs... T) *ContainChecker[T, A] {
	if !a.fits { return a }

	for _, o := range objs {
		if _, ok := a.dict[o]; !ok {
			a.fits = false
			return a
		}
	}

	return a
}

func (a *ContainChecker[T, A]) Not(objs... T) *ContainChecker[T, A] {
	if !a.fits { return a }

	for _, o := range objs {
		if _, ok := a.dict[o]; ok {
			a.fits = false
			return a
		}
	}

	return a
}

func (a *ContainChecker[T, A]) Fits() bool {
	return a.fits
}
