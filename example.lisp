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

// dict_new
(let dict_new (lambda (lambda x 0)))
// dict_get d[x]
(let dict_get (lambda d x (d x)))

// dict_set d[x] = y
(let dict_set (lambda d x y (
    lambda z (
        case z
            x y
            _ (dict_get d z)
    )
)))

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

// dict example
(let d (dict_new))                                      // new dict
(let d (dict_set d 2 300))                                 // set value
(let d (dict_set d 3 500))                                 // set value
(let d (dict_set d 2 200))                                 // set value
(output (dict_get d 2) 11)                                 // should print 200

// end dict example

(output (div 6 2))                                      // test extension

(output 'hello word')
(let data '{
    "name": "khanh",
    "age": 28
}')
(let data (json_add_field data 'subject' 'math'))        // test json extension
(output data)

(let x (input))                                           // waiting for user input
(output (fibonacci x) 11)                                // print the x-th fibonacci

