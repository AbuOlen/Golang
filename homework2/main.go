package main

import (
	"fmt"
)

func main() {
	fmt.Println(FibonacciIterative(11))
	fmt.Println(FibonacciRecursive(11))
	fmt.Println(IsPrime(73))
	fmt.Println(IsBinaryPalindrome(0b1001001))
	fmt.Println(ValidParentheses("[{}}]"))
	fmt.Println(Increment("1001001"))
}

/*
89
89
true
true
false
74
*/
