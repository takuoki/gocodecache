package gocodecache

type initializeOptions struct {
	loadFirstKeys map[string]struct{}
}

type InitializeOption interface {
	apply(*initializeOptions)
}

type funcInitializeOption struct {
	f func(*initializeOptions)
}

func (fdo *funcInitializeOption) apply(do *initializeOptions) {
	fdo.f(do)
}

func newFuncInitializeOption(f func(*initializeOptions)) *funcInitializeOption {
	return &funcInitializeOption{
		f: f,
	}
}

func WithLoadFirstKeys(firstKeys ...string) InitializeOption {
	var m map[string]struct{}
	if len(firstKeys) > 0 {
		m = map[string]struct{}{}
	}
	for _, k := range firstKeys {
		m[k] = struct{}{}
	}
	return newFuncInitializeOption(func(o *initializeOptions) {
		o.loadFirstKeys = m
	})
}

func defaultInitializeOptions() initializeOptions {
	return initializeOptions{
		loadFirstKeys: nil, // nil means all keys
	}
}
