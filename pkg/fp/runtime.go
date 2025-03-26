package fp

import (
	"fmt"
	"strconv"
)

type Runtime interface {
	Eval(*Block) int
}

func NewRuntime() Runtime {
	return &runtime{
		funcImplDict: make(map[string]funcImpl),
		varStack:     []map[string]int{make(map[string]int)},
	}
}

type funcImpl struct {
	paramNameList  []string
	implementation *Block
}
type runtime struct {
	funcImplDict map[string]funcImpl
	varStack     []map[string]int
}

func (r *runtime) Eval(block *Block) int {
	switch block.Type {
	case BLOCKTYPE_NAME:
		// convert to number
		val, err := strconv.Atoi(block.Name)
		if err == nil {
			return val
		}
		// find all variables from top frame to bottom frame
		// NOTE: pure functions will find always find it at the top frame - can detect non-pure function
		for i := len(r.varStack) - 1; i >= 0; i-- {
			if val, ok := r.varStack[i][block.Name]; ok {
				return val
			}
		}
		panic("runtime error")
	case BLOCKTYPE_EXPR:
		switch block.Name {
		case "let":
			return r.builtinLet(block)
		case "func":
			return r.builtinFunc(block)
		case "case":
			return r.builtinCase(block)
		case "input":
			return r.builtinInput(block)
		case "output":
			return r.builtinOutput(block)
		case "sign":
			return r.builtinSign(block)
		case "tail":
			return r.builtinTail(block)
		case "add":
			return r.builtinAdd(block)
		case "sub":
			return r.builtinSub(block)
		default: // user-function application
			f := r.funcImplDict[block.Name]
			// evaluate argument
			localVarDict := map[string]int{}
			for i, arg := range block.Args {
				localVarDict[f.paramNameList[i]] = r.Eval(arg)
			}
			// push new variable stack
			r.varStack = append(r.varStack, localVarDict)
			// evaluate implementation after having argument
			val := r.Eval(f.implementation)
			// pop from variable stack
			r.varStack = r.varStack[:len(r.varStack)-1]
			return val
		}
	default:
		panic("runtime error")
	}
}

// builtinCase : process cases
func (r *runtime) builtinCase(block *Block) int {
	cond := r.Eval(block.Args[0])
	i := func(cond int, args []*Block) int {
		for i := 1; i < len(args); i += 2 {
			arg := args[i]
			// process wildcard independently
			if arg.Type == BLOCKTYPE_NAME && arg.Name == "_" {
				return i
			}
			// process normal case
			if cond == r.Eval(arg) {
				return i
			}
		}
		panic("runtime error")
	}(cond, block.Args)

	return r.Eval(block.Args[i+1])
}

// builtinOutput : evaluate the list of expressions and print
func (r *runtime) builtinOutput(block *Block) int {
	for _, arg := range block.Args {
		fmt.Printf("%d ", r.Eval(arg))
	}
	fmt.Printf("\n")
	return len(block.Args)
}

// builtinLet : evaluate the expression and assign to local variable
func (r *runtime) builtinLet(block *Block) int {
	name := block.Args[0].Name
	value := r.Eval(block.Args[1])
	r.varStack[len(r.varStack)-1][name] = value
	return value
}

// builtinInput : get input from stdin and assign to local variable
func (r *runtime) builtinInput(block *Block) int {
	name := block.Args[0].Name
	var value int
	_, err := fmt.Scan(&value)
	if err != nil {
		panic(err)
	}
	r.varStack[len(r.varStack)-1][name] = value
	return 0
}

// builtinFunc : function definition, save function implementation
func (r *runtime) builtinFunc(block *Block) int {
	name := block.Args[0].Name
	var paramNameList []string
	for i := 1; i < len(block.Args)-1; i++ {
		paramNameList = append(paramNameList, block.Args[i].Name)
	}
	r.funcImplDict[name] = funcImpl{
		paramNameList:  paramNameList,
		implementation: block.Args[len(block.Args)-1],
	}
	return 0
}

// builtinAdd : sum
func (r *runtime) builtinAdd(block *Block) int {
	value := 0
	// evaluate all arguments then return the sum
	// NOTE : if functions are pure, this can be done in parallel
	for _, arg := range block.Args {
		value += r.Eval(arg)
	}
	return value
}

// builtinTail : similar to sum but get the last one
func (r *runtime) builtinTail(block *Block) int {
	value := 0
	// evaluate all arguments then return the last one
	// NOTE : if functions are pure, this can be done in parallel
	for _, arg := range block.Args {
		value = r.Eval(arg)
	}
	return value
}

// builtinSub : subtract
func (r *runtime) builtinSub(block *Block) int {
	// subtraction
	return r.Eval(block.Args[0]) - r.Eval(block.Args[1])
}

// builtinSign : sign function
func (r *runtime) builtinSign(block *Block) int {
	value := r.Eval(block.Args[0])
	switch {
	case value > 0:
		return +1
	case value < 0:
		return -1
	default:
		return 0
	}
}
