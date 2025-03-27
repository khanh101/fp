package fp

type StepOption func(*stepOption) *stepOption
type stepOption struct {
	tailCallOptimization bool
}

type RuntimeOption func(*runtimeOption) *runtimeOption
type runtimeOption struct {
	debug     bool
	parseName func(Name) (interface{}, error)
}

func ParseNameOption(f func(Name) (interface{}, error)) RuntimeOption {
	return func(o *runtimeOption) *runtimeOption {
		o.parseName = f
		return o
	}
}

func DebugOption(debug bool) RuntimeOption {
	return func(o *runtimeOption) *runtimeOption {
		o.debug = debug
		return o
	}
}

func TCOStepOption(tco bool) StepOption {
	return func(o *stepOption) *stepOption {
		o.tailCallOptimization = false // TODO - debug tailcall
		return o
	}
}
