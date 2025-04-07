package fp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	DETECT_NONPURE        = false
	MAX_STACK_DEPTH       = 1000
	TAILCALL_OPTIMIZATION = true
)

func (r *Runtime) searchOnStack(name String) (Object, error) {
	for i := len(r.Stack) - 1; i >= 0; i-- {
		if o, ok := r.Stack[i][name]; ok {
			if DETECT_NONPURE {
				if i != 0 && i < len(r.Stack)-1 {
					_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
				}
			}
			return o, nil
		}
	}
	return nil, fmt.Errorf("object not found %s", name)
}

var InterruptError = errors.New("interrupt")
var TimeoutError = errors.New("timeout")
var StackOverflowError = errors.New("stack overflow")

type stepOptions struct {
	tailCall bool
}

func getOptionsFromContext(ctx context.Context) *stepOptions {
	if o, ok := ctx.Value("step_options").(*stepOptions); ok {
		return o
	}
	// return default option
	return &stepOptions{
		tailCall: false,
	}
}

// Step -
func (r *Runtime) Step(ctx context.Context, expr Expr) (Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	options := getOptionsFromContext(ctx)

	deadline, ok := ctx.Deadline()
	if ok && time.Now().After(deadline) {
		return nil, TimeoutError
	}
	if len(r.Stack) > MAX_STACK_DEPTH {
		return nil, StackOverflowError
	}
	select {
	case <-ctx.Done():
		return nil, InterruptError
	default:
		switch expr := expr.(type) {
		case NameExpr:
			var v Object
			// parse name
			v, err := r.parseLiteral(String(expr))
			if err == nil {
				return v, nil
			}
			// find in stack for variable
			return r.searchOnStack(String(expr))

		case LambdaExpr:
			f, err := r.searchOnStack(String(expr.Name))
			if err != nil {
				return nil, err
			}
			switch f := f.(type) {
			case Lambda:
				// 1. evaluate arguments
				args, err := r.stepMany(ctx, expr.Args...)
				if err != nil {
					return nil, err
				}
				// 2. add argument to local Frame
				localFrame := make(Frame).Update(f.Frame)
				for i := 0; i < len(f.Params); i++ {
					localFrame[f.Params[i]] = args[i]
				}
				// 3. push Frame to Stack
				if options.tailCall {
					r.Stack[len(r.Stack)-1].Update(localFrame)
				} else {
					r.Stack = append(r.Stack, localFrame)
				}
				// 4. exec function
				v, err := r.Step(ctx, f.Impl)
				if err != nil {
					return nil, err
				}
				// 5. pop Frame from Stack
				if !options.tailCall {
					r.Stack = r.Stack[:len(r.Stack)-1]
				}
				return v, nil
			case Module:
				return f.Exec(ctx, r, expr)
			default:
				return nil, fmt.Errorf("function or module %s found but wrong type", expr.Name.String())
			}
		default:
			return nil, fmt.Errorf("runtime error: unknown expression type")
		}
	}
}

func (r *Runtime) stepMany(ctx context.Context, exprList ...Expr) ([]Object, error) {
	var outputs []Object
	if len(exprList) != 0 {
		for i, expr := range exprList {
			if TAILCALL_OPTIMIZATION {
				if i == len(exprList)-1 && len(exprList) >= 2 {
					opts := getOptionsFromContext(ctx)
					opts.tailCall = true
					ctx = context.WithValue(ctx, "step_options", opts)
				}
			}
			v, err := r.Step(ctx, expr)
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, v)
		}

	}
	return outputs, nil
}
