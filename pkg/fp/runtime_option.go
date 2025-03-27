package fp

type StepOption func(*stepOption) *stepOption
type stepOption struct {
	tailCallOptimization bool
}

func defaultStepOption() *stepOption {
	return &stepOption{
		tailCallOptimization: false,
	}
}

func WithTailCallOptimization(o *stepOption) *stepOption {
	o.tailCallOptimization = false // TODO - debug tailcall
	return o
}
