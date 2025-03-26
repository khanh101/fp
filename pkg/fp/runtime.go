package fp

import (
	"fmt"
	"strconv"
)

type Func struct {
	ParamNameList  []string
	Implementation *Block
}
type Runtime struct {
	FuncMap     map[string]Func
	VarMapStack []map[string]int
}

func (r *Runtime) Eval(block *Block) int {
	if block.Type == BLOCKTYPE_LITERAL {
		for i := len(r.VarMapStack) - 1; i >= 0; i-- {
			if val, ok := r.VarMapStack[i][block.Name]; ok {
				return val
			}
		}
		val, err := strconv.Atoi(block.Name)
		if err != nil {
			panic(err)
		}
		return val
	}
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
	default:
		// new frame
		f := r.FuncMap[block.Name]
		varMap := map[string]int{}
		for i := 0; i < len(f.ParamNameList); i++ {
			varMap[f.ParamNameList[i]] = r.Eval(block.Args[i])
		}
		r.VarMapStack = append(r.VarMapStack, varMap)
		val := r.Eval(f.Implementation)
		r.VarMapStack = r.VarMapStack[:len(r.VarMapStack)-1]
		return val
	}
}

func (r *Runtime) builtinSign(block *Block) int {
	val := r.Eval(block.Args[0])
	if val == 0 {
		return 0
	}
	if val > 0 {
		return +1
	}
	if val < 0 {
		return -1
	}
	panic("runtime error")
}

func (r *Runtime) builtinTail(block *Block) int {
	v := 0
	for _, arg := range block.Args {
		v = r.Eval(arg)
	}
	return v
}
func (r *Runtime) builtinAdd(block *Block) int {
	v := 0
	for _, arg := range block.Args {
		v += r.Eval(arg)
	}
	return v
}
func (r *Runtime) builtinSub(block *Block) int {
	return r.Eval(block.Args[0]) - r.Eval(block.Args[1])
}

func (r *Runtime) builtinCase(block *Block) int {
	val := r.Eval(block.Args[0])
	for i := 1; i < len(block.Args); i += 2 {
		if block.Args[i].Type == BLOCKTYPE_LITERAL && block.Args[i].Name == "_" {
			// wildcard
			return r.Eval(block.Args[i+1])
		}
		caseVal := r.Eval(block.Args[i])
		if caseVal == val {
			return r.Eval(block.Args[i+1])
		}
	}
	panic("runtime error")
}

func (r *Runtime) builtinOutput(block *Block) int {
	for _, arg := range block.Args {
		fmt.Printf("%d ", r.Eval(arg))
	}
	fmt.Printf("\n")
	return len(block.Args)
}

func (r *Runtime) builtinLet(block *Block) int {
	name := block.Args[0].Name
	val := r.Eval(block.Args[1])
	r.VarMapStack[len(r.VarMapStack)-1][name] = val
	return 0
}
func (r *Runtime) builtinInput(block *Block) int {
	name := block.Args[0].Name
	var val int
	_, err := fmt.Scan(&val)
	if err != nil {
		panic(err)
	}
	r.VarMapStack[len(r.VarMapStack)-1][name] = val
	return 0
}
func (r *Runtime) builtinFunc(block *Block) int {
	name := block.Args[0].Name
	var paramNameList []string
	for i := 1; i < len(block.Args)-1; i++ {
		paramNameList = append(paramNameList, block.Args[i].Name)
	}
	r.FuncMap[name] = Func{
		ParamNameList:  paramNameList,
		Implementation: block.Args[len(block.Args)-1],
	}
	return 0
}
