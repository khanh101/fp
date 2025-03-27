package fp

type stepOption struct {
	tailCallOptimization bool
}

func WithTailCallOptimization(o *stepOption) *stepOption {
	o.tailCallOptimization = false // TODO - debug tailcall
	return o
}
