```lisp
welcome to fp repl! type <function or module name> for help
>>>add
extension: (add 1 (add 2 3) 3) - exec a sequence of expressions and return the sum
>>>append
extension: (append l 2 (add 1 1)) - append elements into list l and return a new list
>>>case
module: (case x 1 2 4 5) - case, if x=1 then return 3, if x=4 the return 5
>>>del
module: (del x) - delete variable x
>>>div
extension: division
>>>lambda
module: (lambda x y (add x y) - declare a function
>>>let
module: (let x 3) - assign value 3 to local variable x
>>>list
extension: (list 1 2 (lambda x (add x 1))) - make a list
>>>map
extension: (map l (lambda y (add 1 y))) - map
>>>mod
extension: (mod 2 (add 1 1)) - exec two expressions and return modulo
>>>mul
extension: (mul 1 (add 2 3) 3) - exec a sequence of expressions and return the product
>>>peak
extension: (peak l 3 2) - get elem from list (can get multiple elements)
>>>print
extension: (print 1 x (lambda 3)) - print values
>>>reset
module: (reset) - reset stack
>>>sign
extension: (sign 3) - exec an expression and return the sign
>>>slice
extension: (slice l 2 3) - make a slice of a list l[2, 3]
>>>stack
module: (stack) - get stack
>>>sub
extension: (sub 2 (add 1 1)) - exec two expressions and return difference
>>>tail
extension: (tail (print 1) (print 2) 3) - exec a sequence of expressions and return the last one
>>>type
extension: (type x 1 (lambda y (add 1 y))) - get types of objects (can get multiple ones)
>>>unicode
extension: (unicode 72 101 108 108 111 44 32 87 111 114 108 100 33) - convert a list of integers into string - this is just for hello world
```
