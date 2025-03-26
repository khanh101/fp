# FP

A simple functional programming language in 130 lines of code with `functions as first-class citizens`. $\mathbb{F}_p$ also denotes the finite field of order $p$ 


## How to implement data structures like list, dict

list is a vector of integers is $\mathbb{Z}^{\mathbb{N}} \cong \mathbb{Z}$ so we're done. Similarly, for any other data structures

it is also possible to define list as a function $\mathbb{N} \to \mathbb{Z}$, operation on list is composition of function

floating point is useful

## How to handle infix operator?

translate `[<expr_1> <name> <expr_2>]` into `(<name> <expr_1> <expr_2>)` - not yet implemented

## Isn't `(let x 3)` equivalent to `(let x (lambda 3))`?

yes, if functions are pure, then we can consider `(let x <expr>)` as a pure function of the form `(func x <expr>)`. 
however, if functions are not pure, if `x` is defined locally, `(let f (lambda (x + 3)))` and `let f (x + 3)` are different
since variables are evaluated at definition but functions are only evaluated when it is called,
that is if we pass `f` outside of the function, it no longer valid.
in the code below, i gave an example with `(let x_v (output 2 5))` and `(func x_f (output 2 6))`

## How to handle higher-order functions

higher-order function is already implemented

## Performance improvement

if we assume functions are pure, one can consider the whole program as a set of expressions (with some dependencies of `let`)
each function call only need its own variable scope, they can execute every expression at the same time (possibly with some waiting for `let` statement) 

## But can it run Doom?

no 😅

## language specs

- program : a list of expression
- name and expression: name is a string of characters, e.g. `x`, `mul`, and expression is enclosed with parentheses starting with a name, e.g `(let x 3)`, `(add 1 2)`
- evaluation: in run time, name and expression have an associated value
    - name is evaluated using a pool of variables; in code, it is `varStack`. if a name is not of a variable name declared using `let` or `input`, it is undefined behavior
    - expression is evaluated using its name

- builtin functions: `let, func, case, sign, add, sub, tail, input, output, lambda, global`
```
(let <name> <expr>)                                          - assign value of <expr> into local variable <name>
(lambda <name_1> ... <name_n> <expr>)                        - declare an anonymous function
(case <cond> <expr_1> <expr_2>... <key_{n-1}> <expr_n>)      - branching, if <cond> = <key_i> for i odd, return <expr_{i+1}>
(sign <expr>)                                                - return (-1), 0, (+1) according to sign of <expr>
(add <expr_1> ... <expr_n>)                                  - add
(sub <expr_1> <expr_2>)                                      - subtract
(tail <expr_1> ... <expr_n>                                  - evaluate all expressions then return the last one
                                                               (use to declare local variables, do multistep calculation)
(input <name>)                                               - read stdin and assign into <name>
(output <expr_1> ... <expr_n>                                - write to stdout
```

- wildcard symbol: `_` is a special symbol used in `case` to mark every other cases
- no match is `case` is an undefined behavior

## a simple program

```
// define multiplication
(let mul
    (lambda x y
        (case (sign y)                         // mul: (x y) -> xy
            0 0                                             // if y = 0, return 0
            -1 (sub 0 (mul x (sub 0 y)))                    // if y < 0, return 0 - x(-y)
            +1 (add x (mul x (sub y 1)))                    // if y > 0, return x + x(y-1)
        )
    )
)

// define modulo
(let mod
    (lambda x y
        (tail                                  // mul: (x y) -> x % y // defined only for positive y
            (let z (sub x y))                               // local var z = x - y
            (output z x y 6)                                // print local value of z (with label 6)
            (case (sign z)
                +1 (mod z y)                                // if x > y, return (x - y) % y
                0  0                                        // if x = y, return 0
                -1 x                                        // if x < y, return x
            )
        )
    )
)

// define fibonacci

(let fibonacci
    (lambda x
        (case (sign (sub x 1))
            1 (tail
                (let y (fibonacci (sub x 1)))
                (let z (fibonacci (sub x 2)))
                (add y z)
            )
            _ x
        )
    )
)

// partial function using lambda
(let addx
    (lambda x
        (lambda y (add x y))
    )
)

(let z 20)
(output z 1)                                            // print z=20 (with label 1)
(output (mul 13 -17) 2)                                 // print 13 * (-17) (with label 2)
(output (mod 17  13) 3)                                 // print 17 % 13 (with label 3)
(output z 4)                                            // print z=20 again (with label 4), verify that the other z is an actual local variable

(let x_v (output 2 5))                                  // declare x_v - (output 2 5) is executed immediately
(let x_f (lambda (output 2 6)))                         // declare x_f - (output 2 6) is not executed immediately
(output 7)                                              // for debugging
(x_f)                                                   // apply x_f - (output 2 6) is executed

(let f (lambda x (add x 1)))                            // define lambda
(output f)                                              // print lambda
(output (f 21) 8)                                       // print 21 + 1 using lambda

(let t 3)
(let add3 (addx t))                                     // partial function
(output (add3 14) 9)

(input x)                                               // waiting for user input
(output (fibonacci x) 10)                                // print the x-th fibonacci (with label 5)

```
