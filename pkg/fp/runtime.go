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
	case "print":
		return r.builtinPrint(block)
	case "sign":
		return r.builtinSign(block)
	case "add":
		return r.builtinAdd(block)
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
	if block.Name != "sign" {
		panic("runtime error")
	}
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

func (r *Runtime) builtinAdd(block *Block) int {
	if block.Name != "add" {
		panic("runtime error")
	}
	v := 0
	for _, arg := range block.Args {
		v += r.Eval(arg)
	}
	return v
}

func (r *Runtime) builtinCase(block *Block) int {
	if block.Name != "case" {
		panic("runtime error")
	}
	// val
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

func (r *Runtime) builtinPrint(block *Block) int {
	if block.Name != "print" {
		panic("runtime error")
	}
	val := r.Eval(block.Args[0])
	fmt.Print(val)
	return val
}

func (r *Runtime) builtinLet(block *Block) int {
	if block.Name != "let" {
		panic("runtime error")
	}
	name := block.Args[0].Name
	val := r.Eval(block.Args[1])
	r.VarMapStack[len(r.VarMapStack)-1][name] = val
	return 0
}

func (r *Runtime) builtinFunc(block *Block) int {
	if block.Name != "func" {
		panic("runtime error")
	}
	name := block.Args[0].Name
	params := []string{}
	param0 := block.Args[1].Name
	params = append(params, param0)
	for _, param := range block.Args[1].Args {
		params = append(params, param.Name)
	}
	r.FuncMap[name] = Func{
		ParamNameList:  params,
		Implementation: block.Args[2],
	}
	return 0
}
