# FP

A simple functional programming language in 200 lines of code. $\mathbb{F}_p$ also denotes the finite field of order $p$ 


## How to handle higher-order functions

I haven't implemented it yet, but the idea is as follows
- Assign each function to an integer
- builtin `(partial myfunc x y)`
    - evaluate `x`
    - match the integer with a function of 2 parameters, let's say `sub`
    - make new `funcImpl` from `sub`, save evaluation of `y` into `funcImpl`

- whenever `(myfunc z)` is evaluated, apply all partial arguments
- however, one has to handle functions declared locally. So, high-order functions will make the implementation more complicated


## language specs

- name and expression: name is a string of characters, e.g. `x`, `mul`, and expression is enclosed with parentheses starting with a name, e.g `(let x 3)`, `(add 1 2)`
- evaluation: in run time, name and expression have an associated value
    - name is evaluated using a pool of variables; in code, it is `varDictStack`. if a name is not of a variable name declared using `let` or `input`, it is undefined behavior
    - expression is evaluated using its name

- builtin functions: `let, func, case, sign, add, sub, tail, input, output`
```
(let <name> <expr>)                                          - assign value of <expr> into <name>, return 0
(func <name> [<name_1> ... <name_n>] <expr>)                 - declare a function <name> with n parameters, return 0
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
(
    func mul x y (case (sign y)                         // mul: (x y) -> xy
        0 0                                             // if y = 0, return 0
        -1 (sub 0 (mul x (sub 0 y)))                    // if y < 0, return 0 - x(-y)
        +1 (add x (mul x (sub y 1)))                    // if y > 0, return x + x(y-1)
    )
)

// define modulo
(
    func mod x y (tail                                  // mul: (x y) -> x % y // defined only for positive y
        (let z (sub x y))                               // local var z = x - y
        (output z x y 6)                                // print local value of z (with label 6)
        (case (sign z)
            +1 (mod z y)                                // if x > y, return (x - y) % y
            0  0                                        // if x = y, return 0
            -1 x                                        // if x < y, return x
        )
    )
)

// define fibonacci

(
    func fibonacci x (case (sign (sub x 1))
        1 (tail
            (let y (fibonacci (sub x 1)))
            (let z (fibonacci (sub x 2)))
            (add y z)
        )
        _ x
    )
)

(let z 20)
(output z 1)                                            // print z=20 (with label 1)
(output (mul 13 -17) 2)                                 // print 13 * (-17) (with label 2)
(output (mod 17  13) 3)                                 // print 17 % 13 (with label 3)
(output z 4)                                            // print z=20 again (with label 1), verify that the other z is an actual local variable
(input x)                                               // waiting for user input
(output (fibonacci x) 5)                                // print the x-th fibonacci (with label 5)
```
