# FP

A simple functional programming language in 130 lines of code with `functions as first-class citizens`. $\mathbb{F}_p$ also denotes the finite field of order $p$ 

## How to use?

A REPL is available by running `go run bin/repl/main.go`. An example program is at `example.lisp`.

Have fun 🤗

## How to implement data structures like list, dict

data structures can be implemented using functions, example below

floating point is useful

## How to handle infix operator?

translate `[<expr_1> <name_1> <expr_2> <name_2> <expr_3>]` into `(<name_1> <expr_1> (<name_2> <expr_2> <expr_3>))` - todo 
## Isn't `(let x 3)` equivalent to `(let x (lambda 3))`?

yes, if functions are pure, then we can consider `(let x <expr>)` as a pure function of the form `(let x (lambda <expr>))`. 
however, if functions are not pure, if `x` is defined locally, `(let f (lambda (x + 3)))` and `let f (x + 3)` are different
since variables are evaluated at definition but functions are only evaluated when it is called,
that is if we pass `f` outside of the function, it no longer valid.
in the code below, i gave an example with `(let x_v (print 2 5))` and `(func x_f (print 2 6))`

## How to handle higher-order functions

higher-order function is already implemented

## Tail call optimization

WIP

## Performance improvement

if we assume functions are pure, one can consider the whole program as a set of expressions (with some dependencies of `let`)
each function call only need its own variable scope, they can execute every expression at the same time (possibly with some waiting for `let` statement) 

## But can it run Doom?

no 😅

## MANUAL 

- for builtin modules, extensions, see `MANUAL.md`
- wildcard symbol: `_` is a special symbol used in `case` to mark every other cases
- no match is `case` is an undefined behavior

## a simple program

- see `example.lisp`
- hello world ! `(unicode 72 101 108 108 111 44 32 87 111 114 108 100 33)`